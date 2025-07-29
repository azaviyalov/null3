package mood

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/core/auth"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *Handler) {
	jwt := auth.JWTMiddleware()
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
	slog.Debug("GetEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	user, _ := auth.GetUser(c)
	slog.Info("GetEntry request", "id", idParam, "userID", user.ID)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		slog.Warn("invalid id param", "id", idParam)
		return echo.ErrBadRequest.WithInternal(err)
	}
	entry, err := h.service.GetEntry(user.ID, uint(id))
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			slog.Warn("entry not found", "id", id, "userID", user.ID)
			return echo.ErrNotFound.WithInternal(err)
		}
		slog.Error("GetEntry failed", "error", err, "id", id, "userID", user.ID)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) ListEntries(c echo.Context) error {
	slog.Debug("ListEntries handler called", "method", c.Request().Method, "path", c.Path())
	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")
	deletedParam := c.QueryParam("deleted")
	user, _ := auth.GetUser(c)
	slog.Info("ListEntries request params", "limit", limitParam, "offset", offsetParam, "deleted", deletedParam, "userID", user.ID)

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		slog.Warn("invalid limit param", "limit", limitParam)
		limit = 10 // default limit
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset <= 0 {
		slog.Warn("invalid offset param", "offset", offsetParam)
		offset = 0 // default offset
	}

	deleted, err := strconv.ParseBool(deletedParam)
	if err != nil {
		deleted = false // default to not deleted
	}

	entries, err := h.service.ListEntries(user.ID, limit, offset, deleted)
	if err != nil {
		slog.Error("ListEntries failed", "error", err, "userID", user.ID, "limit", limit, "offset", offset, "deleted", deleted)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) CreateEntry(c echo.Context) error {
	slog.Debug("CreateEntry handler called", "method", c.Request().Method, "path", c.Path())
	user, _ := auth.GetUser(c)
	slog.Info("CreateEntry request received", "userID", user.ID)

	var req EditEntryRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("failed to bind CreateEntry request", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Warn("validation failed for CreateEntry", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	resp, err := h.service.CreateEntry(user.ID, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			slog.Warn("invalid entry data", "error", err)
			return echo.ErrBadRequest.WithInternal(err)
		}
		slog.Error("CreateEntry failed", "error", err, "userID", user.ID)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateEntry(c echo.Context) error {
	slog.Debug("UpdateEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	user, _ := auth.GetUser(c)
	slog.Info("UpdateEntry request", "id", idParam, "userID", user.ID)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		slog.Warn("invalid id param", "err", err, "id", idParam, "userID", user.ID)
		return echo.ErrBadRequest.WithInternal(err)
	}

	var req EditEntryRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("failed to bind UpdateEntry request", "error", err, "id", id, "userID", user.ID)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Warn("validation failed for UpdateEntry", "error", err, "id", id, "userID", user.ID)
		return echo.ErrBadRequest.WithInternal(err)
	}

	resp, err := h.service.UpdateEntry(user.ID, uint(id), req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			slog.Warn("invalid entry data", "error", err)
			return echo.ErrBadRequest.WithInternal(err)
		}
		slog.Error("UpdateEntry failed", "error", err, "id", id, "userID", user.ID)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteEntry(c echo.Context) error {
	slog.Debug("DeleteEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	user, _ := auth.GetUser(c)
	slog.Info("DeleteEntry request", "id", idParam, "userID", user.ID)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		slog.Warn("invalid id param", "id", idParam)
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.DeleteEntry(user.ID, uint(id))
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			slog.Warn("entry not found for delete", "id", id, "userID", user.ID)
			return echo.ErrNotFound.WithInternal(err)
		}
		slog.Error("DeleteEntry failed", "error", err, "id", id, "userID", user.ID)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) RestoreEntry(c echo.Context) error {
	slog.Debug("RestoreEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	user, _ := auth.GetUser(c)
	slog.Info("RestoreEntry request", "id", idParam, "userID", user.ID)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		slog.Warn("invalid id param", "id", idParam)
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.RestoreEntry(user.ID, uint(id))
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			slog.Warn("entry not found for restore", "id", id, "userID", user.ID)
			return echo.ErrNotFound.WithInternal(err)
		}
		slog.Error("RestoreEntry failed", "error", err, "id", id, "userID", user.ID)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}
