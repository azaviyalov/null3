package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler, jwtMiddleware echo.MiddlewareFunc) {
	e.POST("/api/auth/login", handler.Login)
	e.POST("/api/auth/logout", handler.Logout, jwtMiddleware)
	e.POST("/api/auth/refresh", handler.Refresh)
	e.GET("/api/auth/me", handler.Me, jwtMiddleware)
}

type Handler struct {
	service        *Service
	validator      *validator.Validate
	config         Config
	stubUserConfig StubUserConfig
}

func NewHandler(service *Service, config Config, stubUserConfig StubUserConfig) *Handler {
	return &Handler{
		service:        service,
		config:         config,
		stubUserConfig: stubUserConfig,
		validator:      validator.New(),
	}
}

func (h *Handler) Login(c echo.Context) error {
	slog.Debug("Login handler called", "method", c.Request().Method, "path", c.Path())
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("failed to bind LoginRequest", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Warn("validation failed for LoginRequest", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, tokenData, err := h.service.Authenticate(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			slog.Warn("invalid credentials", "login", req.Login)
			return echo.ErrUnauthorized.WithInternal(err)
		}
		slog.Error("authentication failed", "error", err, "login", req.Login)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    tokenData.JWT.Value,
		HttpOnly: true,
		Secure:   h.config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokenData.RefreshToken.Value,
		HttpOnly: true,
		Secure:   h.config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.config.RefreshTokenExpiration.Seconds()),
	})
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(c echo.Context) error {
	slog.Debug("Logout handler called", "method", c.Request().Method, "path", c.Path())

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = -1
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Secure = h.config.SecureCookies
	c.SetCookie(cookie)

	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		slog.Warn("failed to get refresh token cookie", "error", err)
	}

	if refreshCookie != nil {
		if err := h.service.InvalidateRefreshToken(refreshCookie.Value); err != nil {
			slog.Error("failed to invalidate refresh token", "error", err)
			return echo.ErrInternalServerError.WithInternal(err)
		}
	}
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   h.config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	return emptyJSON(c, http.StatusOK)
}

func (h *Handler) Me(c echo.Context) error {
	slog.Debug("Me handler called", "method", c.Request().Method, "path", c.Path())
	user, err := GetUser(c)
	if err != nil {
		slog.Warn("failed to get user from context", "error", err)
		return echo.ErrUnauthorized.WithInternal(err)
	}
	meResponse := UserResponse{
		ID: user.ID,
	}
	return c.JSON(http.StatusOK, meResponse)
}

func (h *Handler) Refresh(c echo.Context) error {
	slog.Debug("Refresh handler called", "method", c.Request().Method, "path", c.Path())
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		slog.Warn("failed to get refresh token cookie", "error", err)
		return echo.ErrUnauthorized.WithInternal(ErrRefreshTokenInvalid)
	}
	refreshTokenStr := refreshCookie.Value

	token, err := h.service.GetRefreshToken(refreshTokenStr)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			slog.Warn("refresh token not found", "token", refreshTokenStr)
			return echo.ErrUnauthorized.WithInternal(ErrRefreshTokenInvalid)
		}
		slog.Error("failed to validate refresh token", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	// Check if the refresh token has expired
	if token.ExpiresAt.Before(time.Now()) {
		slog.Warn("refresh token expired", "token", refreshTokenStr)
		return echo.ErrUnauthorized.WithInternal(ErrRefreshTokenInvalid)
	}

	// Invalidate the old refresh token before creating a new one
	if err := h.service.InvalidateRefreshToken(refreshTokenStr); err != nil {
		slog.Error("failed to invalidate old refresh token", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	jwt, err := h.service.GenerateJWT(token.UserID)
	if err != nil {
		slog.Error("failed to generate JWT token", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	newRefreshToken, err := h.service.CreateRefreshToken(token.UserID)
	if err != nil {
		slog.Error("failed to generate new refresh token", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    jwt.Value,
		HttpOnly: true,
		Secure:   h.config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken.Value,
		HttpOnly: true,
		Secure:   h.config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.config.RefreshTokenExpiration.Seconds()),
	})
	return emptyJSON(c, http.StatusOK)
}

func emptyJSON(c echo.Context, status int) error {
	return c.JSON(status, map[string]string{})
}
