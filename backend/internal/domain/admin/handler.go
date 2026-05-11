package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler, adminJWT echo.MiddlewareFunc) {
	e.POST("/api/admin/auth/login", handler.Login)
	e.POST("/api/admin/auth/logout", handler.Logout, adminJWT)
	e.POST("/api/admin/auth/refresh", handler.Refresh)
	e.GET("/api/admin/auth/me", handler.Me, adminJWT)
	e.POST("/api/admin/invites", handler.CreateInvite, adminJWT)
}

type Handler struct {
	accountService *account.Service
	sessionService *session.Service
	validator      *validator.Validate
	config         Config
	sessionConfig  session.Config
}

func NewHandler(accountService *account.Service, sessionService *session.Service, config Config, sessionConfig session.Config) *Handler {
	return &Handler{
		accountService: accountService,
		sessionService: sessionService,
		config:         config,
		sessionConfig:  sessionConfig,
		validator:      validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *Handler) Login(c echo.Context) error {
	var req account.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, tokenData, err := h.accountService.AuthenticateAdmin(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, account.ErrInvalidCredentials) {
			return newHTTPError(http.StatusUnauthorized, "Incorrect admin credentials.", err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	session.SetAdminSessionCookies(c, h.sessionConfig, tokenData)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(c echo.Context) error {
	if refreshCookie, err := c.Cookie(session.AdminRefreshCookieName); err == nil && refreshCookie != nil {
		if err := h.sessionService.InvalidateRefreshToken(c.Request().Context(), refreshCookie.Value); err != nil {
			return echo.ErrInternalServerError.WithInternal(err)
		}
	}

	session.ClearAdminSessionCookies(c, h.sessionConfig)
	return emptyJSON(c, http.StatusOK)
}

func (h *Handler) Me(c echo.Context) error {
	actor, err := session.GetActor(c)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(err)
	}

	user, err := h.accountService.GetUserByID(c.Request().Context(), actor.UserID)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(err)
	}

	return c.JSON(http.StatusOK, account.NewUserResponse(user))
}

func (h *Handler) Refresh(c echo.Context) error {
	refreshCookie, err := c.Cookie(session.AdminRefreshCookieName)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(session.ErrRefreshTokenInvalid)
	}

	res, tokenData, err := h.accountService.RefreshAdminSession(c.Request().Context(), refreshCookie.Value)
	if err != nil {
		if errors.Is(err, session.ErrRefreshTokenInvalid) {
			return echo.ErrUnauthorized.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	session.SetAdminSessionCookies(c, h.sessionConfig, tokenData)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateInvite(c echo.Context) error {
	actor, err := session.GetActor(c)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(err)
	}

	user, err := h.accountService.GetUserByID(c.Request().Context(), actor.UserID)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(err)
	}

	rawToken, invite, err := h.accountService.CreateInvite(c.Request().Context(), user)
	if err != nil {
		if errors.Is(err, account.ErrAdminAccessRequired) {
			return echo.ErrForbidden.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	resp := account.InviteResponse{
		InviteURL: h.frontendURL(c, fmt.Sprintf("/invite/%s", rawToken)),
		ExpiresAt: invite.ExpiresAt,
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) frontendURL(c echo.Context, path string) string {
	if origin := c.Request().Header.Get(echo.HeaderOrigin); origin != "" {
		return strings.TrimRight(origin, "/") + path
	}

	if host := c.Request().Host; host != "" {
		scheme := "http"
		if c.IsTLS() || strings.EqualFold(c.Request().Header.Get("X-Forwarded-Proto"), "https") {
			scheme = "https"
		}
		return fmt.Sprintf("%s://%s%s", scheme, host, path)
	}

	baseURL := strings.TrimRight(h.config.FrontendURL, "/")
	if baseURL == "" {
		baseURL = "http://localhost:4200"
	}
	return baseURL + path
}

func newHTTPError(status int, message string, internal error) error {
	httpError := echo.NewHTTPError(status, message)
	httpError.Internal = internal
	return httpError
}

func emptyJSON(c echo.Context, status int) error {
	return c.JSON(status, struct{}{})
}
