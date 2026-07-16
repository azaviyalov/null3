package journal

import (
	"errors"
	"log/slog"
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
	e.GET("/api/journal/mood/entries", h.ListMoodEntries, jwt)
	e.GET("/api/journal/mood/entries/:id", h.GetMoodEntry, jwt)
	e.POST("/api/journal/mood/entries", h.CreateMoodEntry, jwt)
	e.PUT("/api/journal/mood/entries/:id", h.UpdateMoodEntry, jwt)
	e.DELETE("/api/journal/mood/entries/:id", h.DeleteMoodEntry, jwt)
	e.POST("/api/journal/mood/entries/:id/restore", h.RestoreMoodEntry, jwt)

	e.GET("/api/journal/diary/entries", h.ListDiaryEntries, jwt)
	e.GET("/api/journal/diary/entries/:id", h.GetDiaryEntry, jwt)
	e.POST("/api/journal/diary/entries", h.CreateDiaryEntry, jwt)
	e.PUT("/api/journal/diary/entries/:id", h.UpdateDiaryEntry, jwt)
	e.DELETE("/api/journal/diary/entries/:id", h.DeleteDiaryEntry, jwt)
	e.POST("/api/journal/diary/entries/:id/restore", h.RestoreDiaryEntry, jwt)
}

func (h *Handler) GetMoodEntry(c echo.Context) error {
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	entry, err := h.service.GetMoodEntry(c.Request().Context(), actor.UserID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, NewMoodEntryResponse(entry))
}

func (h *Handler) ListMoodEntries(c echo.Context) error {
	limit, offset, deleted := parsePagination(c)
	actor, _ := session.GetActor(c)
	entries, err := h.service.ListMoodEntries(c.Request().Context(), actor.UserID, limit, offset, deleted)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) CreateMoodEntry(c echo.Context) error {
	actor, _ := session.GetActor(c)
	var req MoodEditEntryRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	entry, err := h.service.CreateMoodEntry(c.Request().Context(), actor.UserID, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusCreated, NewMoodEntryResponse(entry))
}

func (h *Handler) UpdateMoodEntry(c echo.Context) error {
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	var req MoodEditEntryRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	entry, err := h.service.UpdateMoodEntry(c.Request().Context(), actor.UserID, id, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, NewMoodEntryResponse(entry))
}

func (h *Handler) DeleteMoodEntry(c echo.Context) error {
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	entry, err := h.service.DeleteMoodEntry(c.Request().Context(), actor.UserID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, NewMoodEntryResponse(entry))
}

func (h *Handler) RestoreMoodEntry(c echo.Context) error {
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	entry, err := h.service.RestoreMoodEntry(c.Request().Context(), actor.UserID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, NewMoodEntryResponse(entry))
}

func (h *Handler) GetDiaryEntry(c echo.Context) error {
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	entry, err := h.service.GetDiaryEntry(c.Request().Context(), actor.UserID, id)
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
	actor, _ := session.GetActor(c)
	entries, err := h.service.ListDiaryEntries(c.Request().Context(), actor.UserID, limit, offset, deleted)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) CreateDiaryEntry(c echo.Context) error {
	actor, _ := session.GetActor(c)
	var req DiaryEditEntryRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	entry, err := h.service.CreateDiaryEntry(c.Request().Context(), actor.UserID, req)
	if err != nil {
		if errors.Is(err, core.ErrInvalidItem) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusCreated, NewDiaryEntryResponse(entry))
}

func (h *Handler) UpdateDiaryEntry(c echo.Context) error {
	id, actor, err := parseIDAndActor(c)
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

	entry, err := h.service.UpdateDiaryEntry(c.Request().Context(), actor.UserID, id, req)
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
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	entry, err := h.service.DeleteDiaryEntry(c.Request().Context(), actor.UserID, id)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, NewDiaryEntryResponse(entry))
}

func (h *Handler) RestoreDiaryEntry(c echo.Context) error {
	id, actor, err := parseIDAndActor(c)
	if err != nil {
		return err
	}

	entry, err := h.service.RestoreDiaryEntry(c.Request().Context(), actor.UserID, id)
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

func parseIDAndActor(c echo.Context) (uint, *session.Actor, error) {
	idParam := c.Param("id")
	actor, _ := session.GetActor(c)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		slog.Warn(
			"invalid id param",
			"request_id", c.Response().Header().Get("X-Request-Id"),
			"id", idParam,
			"user_id", actor.UserID,
		)
		return 0, nil, echo.ErrBadRequest.WithInternal(err)
	}
	return uint(id), actor, nil
}
