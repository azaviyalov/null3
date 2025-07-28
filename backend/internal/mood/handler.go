package mood

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *Handler) {
	e.GET("/api/mood/entries", h.ListEntries)
	e.GET("/api/mood/entries/:id", h.GetEntry)
	e.POST("/api/mood/entries", h.CreateEntry)
	e.PUT("/api/mood/entries/:id", h.UpdateEntry)
	e.DELETE("/api/mood/entries/:id", h.DeleteEntry)
	e.POST("/api/mood/entries/:id/restore", h.RestoreEntry)
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ListEntries(c echo.Context) error {
	slog.Debug("ListEntries handler called", "method", c.Request().Method, "path", c.Path())
	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")
	deletedParam := c.QueryParam("deleted")
	slog.Info("ListEntries request params", "limit", limitParam, "offset", offsetParam, "deleted", deletedParam)

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

	entries, err := h.service.ListEntries(0, limit, offset, deleted) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		slog.Error("ListEntries failed", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) GetEntry(c echo.Context) error {
	slog.Debug("GetEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	slog.Info("GetEntry request", "id", idParam)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		slog.Warn("invalid id param", "id", idParam)
		return echo.ErrBadRequest.WithInternal(err)
	}
	entry, err := h.service.GetEntry(0, uint(id)) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			slog.Warn("entry not found", "id", id)
			return echo.ErrNotFound.WithInternal(err)
		}
		slog.Error("GetEntry failed", "error", err, "id", id)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) CreateEntry(c echo.Context) error {
	slog.Debug("CreateEntry handler called", "method", c.Request().Method, "path", c.Path())
	slog.Info("CreateEntry request received")
	var req EditEntryRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("failed to bind CreateEntry request", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	resp, err := h.service.CreateEntry(0, req) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			slog.Warn("invalid entry data", "error", err)
			return echo.ErrBadRequest.WithInternal(err)
		}
		slog.Error("CreateEntry failed", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateEntry(c echo.Context) error {
	slog.Debug("UpdateEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	slog.Info("UpdateEntry request", "id", idParam)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		slog.Warn("invalid id param", "id", idParam)
		return echo.ErrBadRequest.WithInternal(err)
	}

	var req EditEntryRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("failed to bind UpdateEntry request", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}

	resp, err := h.service.UpdateEntry(0, uint(id), req) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			slog.Warn("invalid entry data", "error", err)
			return echo.ErrBadRequest.WithInternal(err)
		}
		slog.Error("UpdateEntry failed", "error", err, "id", id)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteEntry(c echo.Context) error {
	slog.Debug("DeleteEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	slog.Info("DeleteEntry request", "id", idParam)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		slog.Warn("invalid id param", "id", idParam)
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.DeleteEntry(0, uint(id)) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			slog.Warn("entry not found for delete", "id", id)
			return echo.ErrNotFound.WithInternal(err)
		}
		slog.Error("DeleteEntry failed", "error", err, "id", id)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) RestoreEntry(c echo.Context) error {
	slog.Debug("RestoreEntry handler called", "method", c.Request().Method, "path", c.Path())
	idParam := c.Param("id")
	slog.Info("RestoreEntry request", "id", idParam)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		slog.Warn("invalid id param", "id", idParam)
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.RestoreEntry(0, uint(id)) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			slog.Warn("entry not found for restore", "id", id)
			return echo.ErrNotFound.WithInternal(err)
		}
		slog.Error("RestoreEntry failed", "error", err, "id", id)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}
