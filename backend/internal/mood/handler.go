package mood

import (
	"errors"
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
	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")
	deletedParam := c.QueryParam("deleted")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10 // default limit
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset <= 0 {
		offset = 0 // default offset
	}

	deleted, err := strconv.ParseBool(deletedParam)
	if err != nil {
		deleted = false // default to not deleted
	}

	entries, err := h.service.ListEntries(0, limit, offset, deleted) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) GetEntry(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	entry, err := h.service.GetEntry(0, uint(id)) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) CreateEntry(c echo.Context) error {
	var req EditEntryRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	resp, err := h.service.CreateEntry(0, req) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateEntry(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	var req EditEntryRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	resp, err := h.service.UpdateEntry(0, uint(id), req) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteEntry(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.DeleteEntry(0, uint(id)) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) RestoreEntry(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.RestoreEntry(0, uint(id)) // Assuming userID is 0 for simplicity, replace with actual user ID logic
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entry)
}
