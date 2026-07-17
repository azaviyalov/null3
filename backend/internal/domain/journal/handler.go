package journal

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func RegisterRoutes(e *echo.Echo, h *Handler, jwt echo.MiddlewareFunc) {
	e.GET("/api/journal/mood-records", h.ListMoodRecords, jwt)
	e.GET("/api/journal/mood-records/:id", h.GetMoodRecord, jwt)
	e.POST("/api/journal/mood-records", h.CreateMoodRecord, jwt)
	e.PUT("/api/journal/mood-records/:id", h.UpdateMoodRecord, jwt)
	e.DELETE("/api/journal/mood-records/:id", h.DeleteMoodRecord, jwt)
	e.POST("/api/journal/mood-records/:id/restore", h.RestoreMoodRecord, jwt)

	e.GET("/api/journal/diary-entries", h.ListDiaryEntries, jwt)
	e.GET("/api/journal/diary-entries/:id", h.GetDiaryEntry, jwt)
	e.POST("/api/journal/diary-entries", h.CreateDiaryEntry, jwt)
	e.PUT("/api/journal/diary-entries/:id", h.UpdateDiaryEntry, jwt)
	e.DELETE("/api/journal/diary-entries/:id", h.DeleteDiaryEntry, jwt)
	e.POST("/api/journal/diary-entries/:id/restore", h.RestoreDiaryEntry, jwt)
}

func (h *Handler) GetMoodRecord(c echo.Context) error {
	id, userID, err := parseIDAndUserID(c)
	if err != nil {
		return err
	}

	entry, err := h.service.GetMoodRecord(c.Request().Context(), userID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, NewMoodRecordResponse(entry))
}

func (h *Handler) ListMoodRecords(c echo.Context) error {
	limit, offset, deleted := parsePagination(c)
	userID := session.GetUserID(c)
	entries, err := h.service.ListMoodRecords(c.Request().Context(), userID, limit, offset, deleted)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) CreateMoodRecord(c echo.Context) error {
	userID := session.GetUserID(c)
	var req MoodEditRecordRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	entry, err := h.service.CreateMoodRecord(c.Request().Context(), userID, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusCreated, NewMoodRecordResponse(entry))
}

func (h *Handler) UpdateMoodRecord(c echo.Context) error {
	id, userID, err := parseIDAndUserID(c)
	if err != nil {
		return err
	}

	var req MoodEditRecordRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.UpdateMoodRecord(c.Request().Context(), userID, id, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, NewMoodRecordResponse(entry))
}

func (h *Handler) DeleteMoodRecord(c echo.Context) error {
	id, userID, err := parseIDAndUserID(c)
	if err != nil {
		return err
	}

	entry, err := h.service.DeleteMoodRecord(c.Request().Context(), userID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, NewMoodRecordResponse(entry))
}

func (h *Handler) RestoreMoodRecord(c echo.Context) error {
	id, userID, err := parseIDAndUserID(c)
	if err != nil {
		return err
	}

	entry, err := h.service.RestoreMoodRecord(c.Request().Context(), userID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, NewMoodRecordResponse(entry))
}

func (h *Handler) GetDiaryEntry(c echo.Context) error {
	id, userID, err := parseIDAndUserID(c)
	if err != nil {
		return err
	}

	entry, err := h.service.GetDiaryEntry(c.Request().Context(), userID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, NewDiaryEntryResponse(entry))
}

func (h *Handler) ListDiaryEntries(c echo.Context) error {
	limit, offset, deleted := parsePagination(c)
	userID := session.GetUserID(c)
	entries, err := h.service.ListDiaryEntries(c.Request().Context(), userID, limit, offset, deleted)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) CreateDiaryEntry(c echo.Context) error {
	userID := session.GetUserID(c)
	var req DiaryEditEntryRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	entry, err := h.service.CreateDiaryEntry(c.Request().Context(), userID, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusCreated, NewDiaryEntryResponse(entry))
}

func (h *Handler) UpdateDiaryEntry(c echo.Context) error {
	id, userID, err := parseIDAndUserID(c)
	if err != nil {
		return err
	}

	var req DiaryEditEntryRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.UpdateDiaryEntry(c.Request().Context(), userID, id, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, NewDiaryEntryResponse(entry))
}

func (h *Handler) DeleteDiaryEntry(c echo.Context) error {
	id, userID, err := parseIDAndUserID(c)
	if err != nil {
		return err
	}

	entry, err := h.service.DeleteDiaryEntry(c.Request().Context(), userID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, NewDiaryEntryResponse(entry))
}

func (h *Handler) RestoreDiaryEntry(c echo.Context) error {
	id, userID, err := parseIDAndUserID(c)
	if err != nil {
		return err
	}

	entry, err := h.service.RestoreDiaryEntry(c.Request().Context(), userID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		if errors.Is(err, core.ErrInvalidItem) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, NewDiaryEntryResponse(entry))
}

func parsePagination(c echo.Context) (int, int, bool) {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	deleted, err := strconv.ParseBool(c.QueryParam("deleted"))
	if err != nil {
		deleted = false
	}
	return limit, offset, deleted
}

func parseIDAndUserID(c echo.Context) (uint, uint, error) {
	idParam := c.Param("id")
	userID := session.GetUserID(c)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, 0, echo.ErrBadRequest.WithInternal(err)
	}
	return uint(id), userID, nil
}
