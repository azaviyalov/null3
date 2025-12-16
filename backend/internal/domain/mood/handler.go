package mood

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/core/auth"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *Handler, jwt echo.MiddlewareFunc) {
	e.GET("/api/mood/entries", h.ListEntries, jwt)
	e.GET("/api/mood/entries/:id", h.GetEntry, jwt)
	e.POST("/api/mood/entries", h.CreateEntry, jwt)
	e.PUT("/api/mood/entries/:id", h.UpdateEntry, jwt)
	e.DELETE("/api/mood/entries/:id", h.DeleteEntry, jwt)
	e.POST("/api/mood/entries/:id/restore", h.RestoreEntry, jwt)
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
	idParam := c.Param("id")
	user, _ := auth.GetUser(c)
	logging.Info(c, "GetEntry request", "id", idParam, "user_id", user.ID)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logging.Warn(c, "invalid id param", "id", idParam)
		return echo.ErrBadRequest.WithInternal(err)
	}
	entry, err := h.service.GetEntry(c.Request().Context(), user.ID, uint(id))
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			logging.Info(c, "entry not found", "id", id, "user_id", user.ID)
			return echo.ErrNotFound.WithInternal(err)
		}
		logging.Error(c, "GetEntry failed", "error", err, "id", id, "user_id", user.ID)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) ListEntries(c echo.Context) error {
	logging.Debug(c, "ListEntries handler called", "method", c.Request().Method, "path", c.Path())
	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")
	deletedParam := c.QueryParam("deleted")
	user, _ := auth.GetUser(c)
	logging.Info(c, "ListEntries request params", "limit", limitParam, "offset", offsetParam, "deleted", deletedParam, "user_id", user.ID)

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		logging.Info(c, "invalid limit param, using default", "limit", limitParam, "default", 10)
		limit = 10 // default limit
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		logging.Info(c, "invalid offset param, using default", "offset", offsetParam, "default", 0)
		offset = 0 // default offset
	}

	deleted, err := strconv.ParseBool(deletedParam)
	if err != nil {
		deleted = false // default to not deleted
	}

	entries, err := h.service.ListEntries(c.Request().Context(), user.ID, limit, offset, deleted)
	if err != nil {
		logging.Error(c, "ListEntries failed", "error", err, "user_id", user.ID, "limit", limit, "offset", offset, "deleted", deleted)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) CreateEntry(c echo.Context) error {
	logging.Debug(c, "CreateEntry handler called", "method", c.Request().Method, "path", c.Path())
	user, _ := auth.GetUser(c)
	logging.Info(c, "CreateEntry request received", "user_id", user.ID)

	var req EditEntryRequest
	if err := c.Bind(&req); err != nil {
		logging.Info(c, "failed to bind CreateEntry request", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		logging.Info(c, "validation failed for CreateEntry", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	resp, err := h.service.CreateEntry(c.Request().Context(), user.ID, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			logging.Info(c, "invalid entry data", "error", err)
			return echo.ErrBadRequest.WithInternal(err)
		}
		logging.Error(c, "CreateEntry failed", "error", err, "user_id", user.ID)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateEntry(c echo.Context) error {
	logging.Debug(c, "UpdateEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	user, _ := auth.GetUser(c)
	logging.Info(c, "UpdateEntry request", "id", idParam, "user_id", user.ID)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logging.Warn(c, "invalid id param", "err", err, "id", idParam, "user_id", user.ID)
		return echo.ErrBadRequest.WithInternal(err)
	}

	var req EditEntryRequest
	if err := c.Bind(&req); err != nil {
		logging.Info(c, "failed to bind UpdateEntry request", "error", err, "id", id, "user_id", user.ID)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		logging.Info(c, "validation failed for UpdateEntry", "error", err, "id", id, "user_id", user.ID)
		return echo.ErrBadRequest.WithInternal(err)
	}

	resp, err := h.service.UpdateEntry(c.Request().Context(), user.ID, uint(id), req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			logging.Info(c, "invalid entry data", "error", err)
			return echo.ErrBadRequest.WithInternal(err)
		}
		logging.Error(c, "UpdateEntry failed", "error", err, "id", id, "user_id", user.ID)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteEntry(c echo.Context) error {
	logging.Debug(c, "DeleteEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	user, _ := auth.GetUser(c)
	logging.Info(c, "DeleteEntry request", "id", idParam, "user_id", user.ID)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logging.Warn(c, "invalid id param", "id", idParam)
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.DeleteEntry(c.Request().Context(), user.ID, uint(id))
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			logging.Info(c, "entry not found for delete", "id", id, "user_id", user.ID)
			return echo.ErrNotFound.WithInternal(err)
		}
		logging.Error(c, "DeleteEntry failed", "error", err, "id", id, "user_id", user.ID)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) RestoreEntry(c echo.Context) error {
	logging.Debug(c, "RestoreEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	user, _ := auth.GetUser(c)
	logging.Info(c, "RestoreEntry request", "id", idParam, "user_id", user.ID)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logging.Warn(c, "invalid id param", "id", idParam)
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.RestoreEntry(c.Request().Context(), user.ID, uint(id))
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			logging.Info(c, "entry not found for restore", "id", id, "user_id", user.ID)
			return echo.ErrNotFound.WithInternal(err)
		}
		logging.Error(c, "RestoreEntry failed", "error", err, "id", id, "user_id", user.ID)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}
