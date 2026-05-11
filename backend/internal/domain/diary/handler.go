package diary

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *Handler, jwt echo.MiddlewareFunc) {
	e.GET("/api/diary/entries", h.ListEntries, jwt)
	e.GET("/api/diary/entries/:id", h.GetEntry, jwt)
	e.POST("/api/diary/entries", h.CreateEntry, jwt)
	e.PUT("/api/diary/entries/:id", h.UpdateEntry, jwt)
	e.DELETE("/api/diary/entries/:id", h.DeleteEntry, jwt)
	e.POST("/api/diary/entries/:id/restore", h.RestoreEntry, jwt)
}

type Handler struct {
	service   *Service
	validator *validator.Validate
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *Handler) GetEntry(c echo.Context) error {
	logging.Debug(c, "GetEntry handler called", "method", c.Request().Method, "path", c.Path())
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	entry, err := h.service.GetEntry(c.Request().Context(), actor.UserID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			logging.Info(c, "diary entry not found", "id", id, "user_id", actor.UserID)
			return echo.ErrNotFound.WithInternal(err)
		}

		logging.Error(c, "GetEntry failed", "error", err, "id", id, "user_id", actor.UserID)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return h.respondWithEntry(c, http.StatusOK, actor.UserID, entry)
}

func (h *Handler) ListEntries(c echo.Context) error {
	logging.Debug(c, "ListEntries handler called", "method", c.Request().Method, "path", c.Path())

	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")
	deletedParam := c.QueryParam("deleted")
	actor, _ := session.GetActor(c)

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	deleted, err := strconv.ParseBool(deletedParam)
	if err != nil {
		deleted = false
	}

	entries, err := h.service.ListEntries(c.Request().Context(), actor.UserID, limit, offset, deleted)
	if err != nil {
		logging.Error(c, "ListEntries failed", "error", err, "user_id", actor.UserID, "limit", limit, "offset", offset, "deleted", deleted)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) CreateEntry(c echo.Context) error {
	logging.Debug(c, "CreateEntry handler called", "method", c.Request().Method, "path", c.Path())
	actor, _ := session.GetActor(c)

	var req EditEntryRequest
	if err := c.Bind(&req); err != nil {
		logging.Info(c, "failed to bind CreateEntry request", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		logging.Info(c, "validation failed for CreateEntry", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.CreateEntry(c.Request().Context(), actor.UserID, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			logging.Info(c, "invalid diary entry data", "error", err, "user_id", actor.UserID)
			return echo.ErrBadRequest.WithInternal(err)
		}

		logging.Error(c, "CreateEntry failed", "error", err, "user_id", actor.UserID)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return h.respondWithEntry(c, http.StatusCreated, actor.UserID, entry)
}

func (h *Handler) UpdateEntry(c echo.Context) error {
	logging.Debug(c, "UpdateEntry handler called", "method", c.Request().Method, "path", c.Path())
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	var req EditEntryRequest
	if err := c.Bind(&req); err != nil {
		logging.Info(c, "failed to bind UpdateEntry request", "error", err, "id", id, "user_id", actor.UserID)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		logging.Info(c, "validation failed for UpdateEntry", "error", err, "id", id, "user_id", actor.UserID)
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.UpdateEntry(c.Request().Context(), actor.UserID, id, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			logging.Info(c, "invalid diary entry update", "error", err, "id", id, "user_id", actor.UserID)
			return echo.ErrBadRequest.WithInternal(err)
		}
		if errors.Is(err, core.ErrItemNotFound) {
			logging.Info(c, "diary entry not found for update", "id", id, "user_id", actor.UserID)
			return echo.ErrNotFound.WithInternal(err)
		}

		logging.Error(c, "UpdateEntry failed", "error", err, "id", id, "user_id", actor.UserID)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return h.respondWithEntry(c, http.StatusOK, actor.UserID, entry)
}

func (h *Handler) DeleteEntry(c echo.Context) error {
	logging.Debug(c, "DeleteEntry handler called", "method", c.Request().Method, "path", c.Path())
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	entry, err := h.service.DeleteEntry(c.Request().Context(), actor.UserID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			logging.Info(c, "diary entry not found for delete", "id", id, "user_id", actor.UserID)
			return echo.ErrNotFound.WithInternal(err)
		}

		logging.Error(c, "DeleteEntry failed", "error", err, "id", id, "user_id", actor.UserID)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return h.respondWithEntry(c, http.StatusOK, actor.UserID, entry)
}

func (h *Handler) RestoreEntry(c echo.Context) error {
	logging.Debug(c, "RestoreEntry handler called", "method", c.Request().Method, "path", c.Path())
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	entry, err := h.service.RestoreEntry(c.Request().Context(), actor.UserID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			logging.Info(c, "diary entry not found for restore", "id", id, "user_id", actor.UserID)
			return echo.ErrNotFound.WithInternal(err)
		}
		if errors.Is(err, core.ErrInvalidItem) {
			logging.Info(c, "invalid diary entry restore", "error", err, "id", id, "user_id", actor.UserID)
			return echo.ErrBadRequest.WithInternal(err)
		}

		logging.Error(c, "RestoreEntry failed", "error", err, "id", id, "user_id", actor.UserID)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return h.respondWithEntry(c, http.StatusOK, actor.UserID, entry)
}

func (h *Handler) respondWithEntry(c echo.Context, statusCode int, userID uint, entry *Entry) error {
	referencedMoodEntries, err := h.service.ListReferencedMoodEntries(c.Request().Context(), userID, entry.ID)
	if err != nil {
		logging.Error(c, "failed to list referenced mood entries", "error", err, "id", entry.ID, "user_id", userID)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(statusCode, NewEntryResponse(entry, referencedMoodEntries))
}

func parseIDAndActor(c echo.Context) (uint, *session.Actor, error) {
	idParam := c.Param("id")
	actor, _ := session.GetActor(c)

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logging.Warn(c, "invalid id param", "id", idParam, "user_id", actor.UserID)
		return 0, nil, echo.ErrBadRequest.WithInternal(err)
	}

	return uint(id), actor, nil
}
