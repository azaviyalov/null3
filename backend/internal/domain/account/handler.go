package account

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler, userJWT echo.MiddlewareFunc) {
	e.POST("/api/auth/login", handler.Login)
	e.POST("/api/auth/logout", handler.Logout, userJWT)
	e.POST("/api/auth/refresh", handler.Refresh)
	e.GET("/api/auth/me", handler.Me, userJWT)
	e.POST("/api/auth/forgot-password", handler.ForgotPassword)
	e.POST("/api/auth/reset-password", handler.ResetPassword)
	e.GET("/api/auth/invites/:token", handler.GetInvite)
	e.POST("/api/auth/invites/:token/register", handler.RegisterWithInvite)
}

type Handler struct {
	service        *Service
	sessionService *session.Service
	config         Config
	sessionConfig  session.Config
}

func NewHandler(service *Service, sessionService *session.Service, config Config, sessionConfig session.Config) *Handler {
	return &Handler{
		service:        service,
		sessionService: sessionService,
		config:         config,
		sessionConfig:  sessionConfig,
	}
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, tokenData, err := h.service.AuthenticateUser(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return newHTTPError(http.StatusUnauthorized, "Incorrect login credentials.", err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	session.SetUserSessionCookies(c, h.sessionConfig, tokenData)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(c echo.Context) error {
	if refreshCookie, err := c.Cookie(session.UserRefreshCookieName); err == nil && refreshCookie != nil {
		if err := h.sessionService.InvalidateRefreshToken(c.Request().Context(), refreshCookie.Value); err != nil {
			return echo.ErrInternalServerError.WithInternal(err)
		}
	}

	session.ClearUserSessionCookies(c, h.sessionConfig)
	return emptyJSON(c, http.StatusOK)
}

func (h *Handler) Me(c echo.Context) error {
	userID := session.GetUserID(c)
	user, err := h.service.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(err)
	}

	return c.JSON(http.StatusOK, NewUserResponse(user))
}

func (h *Handler) Refresh(c echo.Context) error {
	refreshCookie, err := c.Cookie(session.UserRefreshCookieName)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(session.ErrRefreshTokenInvalid)
	}

	res, tokenData, err := h.service.RefreshUserSession(c.Request().Context(), refreshCookie.Value)
	if err != nil {
		if errors.Is(err, session.ErrRefreshTokenInvalid) {
			return echo.ErrUnauthorized.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	session.SetUserSessionCookies(c, h.sessionConfig, tokenData)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetInvite(c echo.Context) error {
	invite, err := h.service.ValidateInvite(c.Request().Context(), c.Param("token"))
	if err != nil {
		if isInviteError(err) {
			return newHTTPError(http.StatusBadRequest, inviteErrorMessage(err), err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, InviteValidationResponse{
		ExpiresAt: invite.ExpiresAt,
	})
}

func (h *Handler) RegisterWithInvite(c echo.Context) error {
	var req InviteRegistrationRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, tokenData, err := h.service.RegisterWithInvite(c.Request().Context(), c.Param("token"), req)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrInvalidItem):
			return newHTTPError(http.StatusBadRequest, clientErrorMessage(err), err)
		case errors.Is(err, ErrInviteInvalid), errors.Is(err, ErrInviteExpired), errors.Is(err, ErrInviteAlreadyUsed):
			return newHTTPError(http.StatusBadRequest, inviteErrorMessage(err), err)
		case errors.Is(err, ErrLoginAlreadyTaken), errors.Is(err, ErrEmailAlreadyTaken):
			return newHTTPError(http.StatusConflict, clientErrorMessage(err), err)
		default:
			return echo.ErrInternalServerError.WithInternal(err)
		}
	}

	session.SetUserSessionCookies(c, h.sessionConfig, tokenData)
	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) ForgotPassword(c echo.Context) error {
	var req ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	rawToken, err := h.service.RequestPasswordReset(c.Request().Context(), req)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	resp := ForgotPasswordResponse{
		Message: "If an account exists for that email, a reset link has been generated.",
	}
	if rawToken != "" {
		resp.ResetURL = h.frontendURL(fmt.Sprintf("/reset-password?token=%s", rawToken))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) ResetPassword(c echo.Context) error {
	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.service.ResetPassword(c.Request().Context(), req); err != nil {
		switch {
		case errors.Is(err, core.ErrInvalidItem):
			return newHTTPError(http.StatusBadRequest, clientErrorMessage(err), err)
		case errors.Is(err, ErrPasswordResetTokenInvalid), errors.Is(err, ErrPasswordResetTokenExpired):
			return newHTTPError(http.StatusBadRequest, resetPasswordErrorMessage(err), err)
		default:
			return echo.ErrInternalServerError.WithInternal(err)
		}
	}

	return c.JSON(http.StatusOK, MessageResponse{
		Message: "Password has been reset successfully.",
	})
}

func (h *Handler) frontendURL(path string) string {
	baseURL := strings.TrimRight(h.config.FrontendURL, "/")
	if baseURL == "" {
		baseURL = "http://localhost:4200"
	}
	return baseURL + path
}

func isInviteError(err error) bool {
	return errors.Is(err, ErrInviteInvalid) || errors.Is(err, ErrInviteExpired) || errors.Is(err, ErrInviteAlreadyUsed)
}

func inviteErrorMessage(err error) string {
	switch {
	case errors.Is(err, ErrInviteAlreadyUsed):
		return "This invite link has already been used."
	case errors.Is(err, ErrInviteExpired):
		return "This invite link has expired."
	default:
		return "This invite link is invalid."
	}
}

func resetPasswordErrorMessage(err error) string {
	if errors.Is(err, ErrPasswordResetTokenExpired) {
		return "This password reset link has expired."
	}
	return "This password reset link is invalid."
}

func clientErrorMessage(err error) string {
	switch {
	case errors.Is(err, ErrLoginAlreadyTaken):
		return "That login is already in use."
	case errors.Is(err, ErrEmailAlreadyTaken):
		return "That email is already in use."
	case strings.HasPrefix(err.Error(), core.ErrInvalidItem.Error()+": "):
		return strings.TrimPrefix(err.Error(), core.ErrInvalidItem.Error()+": ")
	default:
		return "Bad Request"
	}
}

func newHTTPError(status int, message string, internal error) error {
	httpError := echo.NewHTTPError(status, message)
	httpError.Internal = internal
	return httpError
}

func emptyJSON(c echo.Context, status int) error {
	return c.JSON(status, struct{}{})
}
