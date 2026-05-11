package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler, userJWT, adminJWT echo.MiddlewareFunc) {
	e.POST("/api/auth/login", handler.Login)
	e.POST("/api/auth/logout", handler.Logout, userJWT)
	e.POST("/api/auth/refresh", handler.Refresh)
	e.GET("/api/auth/me", handler.Me, userJWT)
	e.POST("/api/auth/forgot-password", handler.ForgotPassword)
	e.POST("/api/auth/reset-password", handler.ResetPassword)
	e.GET("/api/auth/invites/:token", handler.GetInvite)
	e.POST("/api/auth/invites/:token/register", handler.RegisterWithInvite)

	e.POST("/api/admin/auth/login", handler.AdminLogin)
	e.POST("/api/admin/auth/logout", handler.AdminLogout, adminJWT)
	e.POST("/api/admin/auth/refresh", handler.AdminRefresh)
	e.GET("/api/admin/auth/me", handler.AdminMe, adminJWT)
	e.POST("/api/admin/invites", handler.CreateInvite, adminJWT)
}

type Handler struct {
	service   *Service
	validator *validator.Validate
	config    Config
}

func NewHandler(service *Service, config Config) *Handler {
	return &Handler{
		service:   service,
		config:    config,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *Handler) Login(c echo.Context) error {
	return h.login(c, false)
}

func (h *Handler) AdminLogin(c echo.Context) error {
	return h.login(c, true)
}

func (h *Handler) login(c echo.Context, admin bool) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	var (
		res       *UserResponse
		tokenData *UserTokenData
		err       error
	)
	if admin {
		res, tokenData, err = h.service.AuthenticateAdmin(c.Request().Context(), req)
	} else {
		res, tokenData, err = h.service.AuthenticateUser(c.Request().Context(), req)
	}
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			message := "Incorrect login credentials."
			if admin {
				message = "Incorrect admin credentials."
			}
			return newHTTPError(http.StatusUnauthorized, message, err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	scope := userSessionCookies
	if admin {
		scope = adminSessionCookies
	}
	setSessionCookies(c, h.config, scope, tokenData)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(c echo.Context) error {
	return h.logout(c, userSessionCookies)
}

func (h *Handler) AdminLogout(c echo.Context) error {
	return h.logout(c, adminSessionCookies)
}

func (h *Handler) logout(c echo.Context, scope sessionCookies) error {
	if refreshCookie, err := c.Cookie(scope.refreshCookieName); err == nil && refreshCookie != nil {
		if err := h.service.InvalidateRefreshToken(c.Request().Context(), refreshCookie.Value); err != nil {
			return echo.ErrInternalServerError.WithInternal(err)
		}
	}
	clearSessionCookies(c, h.config, scope)
	return emptyJSON(c, http.StatusOK)
}

func (h *Handler) Me(c echo.Context) error {
	return h.me(c)
}

func (h *Handler) AdminMe(c echo.Context) error {
	return h.me(c)
}

func (h *Handler) me(c echo.Context) error {
	user, err := GetUser(c)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(err)
	}
	return c.JSON(http.StatusOK, NewUserResponse(user))
}

func (h *Handler) Refresh(c echo.Context) error {
	return h.refresh(c, userSessionCookies, false)
}

func (h *Handler) AdminRefresh(c echo.Context) error {
	return h.refresh(c, adminSessionCookies, true)
}

func (h *Handler) refresh(c echo.Context, scope sessionCookies, admin bool) error {
	refreshCookie, err := c.Cookie(scope.refreshCookieName)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(ErrRefreshTokenInvalid)
	}

	var (
		res       *UserResponse
		tokenData *UserTokenData
	)
	if admin {
		res, tokenData, err = h.service.RefreshAdminSession(c.Request().Context(), refreshCookie.Value)
	} else {
		res, tokenData, err = h.service.RefreshUserSession(c.Request().Context(), refreshCookie.Value)
	}
	if err != nil {
		if errors.Is(err, ErrRefreshTokenInvalid) {
			return echo.ErrUnauthorized.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	setSessionCookies(c, h.config, scope, tokenData)
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateInvite(c echo.Context) error {
	user, err := GetUser(c)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(err)
	}

	rawToken, invite, err := h.service.CreateInvite(c.Request().Context(), user)
	if err != nil {
		if errors.Is(err, ErrAdminAccessRequired) {
			return echo.ErrForbidden.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}

	resp := InviteResponse{
		InviteURL: h.frontendURL(c, fmt.Sprintf("/invite/%s", rawToken)),
		ExpiresAt: invite.ExpiresAt,
	}
	return c.JSON(http.StatusCreated, resp)
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
	if err := h.validator.Struct(req); err != nil {
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

	setSessionCookies(c, h.config, userSessionCookies, tokenData)
	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) ForgotPassword(c echo.Context) error {
	var req ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
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
		resp.ResetURL = h.frontendURL(c, fmt.Sprintf("/reset-password?token=%s", rawToken))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) ResetPassword(c echo.Context) error {
	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
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
