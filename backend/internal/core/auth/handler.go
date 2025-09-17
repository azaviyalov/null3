package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler, jwtMiddleware echo.MiddlewareFunc) {
	// Regular user authentication routes
	e.POST("/api/auth/login", handler.Login)
	e.POST("/api/auth/logout", handler.Logout, jwtMiddleware)
	e.POST("/api/auth/refresh", handler.Refresh)
	e.GET("/api/auth/me", handler.Me, jwtMiddleware)

	// Admin authentication routes
	e.POST("/api/admin/login", handler.AdminLogin)
	e.POST("/api/admin/logout", handler.AdminLogout, handler.RequireAdminJWT)

	// Admin user management routes
	adminGroup := e.Group("/api/admin", handler.RequireAdminJWT)
	adminGroup.GET("/users", handler.GetUsers)
	adminGroup.POST("/users", handler.CreateUser)
	adminGroup.GET("/users/:id", handler.GetUser)
	adminGroup.PUT("/users/:id", handler.UpdateUser)
	adminGroup.DELETE("/users/:id", handler.DeleteUser)

	// Admin refresh token management routes
	adminGroup.GET("/refresh-tokens", handler.GetRefreshTokens)
	adminGroup.DELETE("/refresh-tokens/:value", handler.DeleteRefreshToken)
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
	return c.JSON(status, struct{}{})
}

// RequireAdminJWT middleware to check for admin JWT
func (h *Handler) RequireAdminJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			slog.Warn("failed to get JWT cookie", "error", err)
			return echo.ErrUnauthorized.WithInternal(err)
		}

		jwt, err := h.service.ParseJWT(cookie.Value)
		if err != nil {
			slog.Warn("failed to parse JWT", "error", err)
			return echo.ErrUnauthorized.WithInternal(err)
		}

		if !h.service.IsAdminJWT(cookie.Value) {
			slog.Warn("non-admin JWT provided for admin endpoint")
			return echo.ErrForbidden
		}

		c.Set("user", &User{ID: jwt.UserID})
		return next(c)
	}
}

func (h *Handler) AdminLogin(c echo.Context) error {
	slog.Debug("AdminLogin handler called", "method", c.Request().Method, "path", c.Path())
	var req AdminLoginRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("failed to bind AdminLoginRequest", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Warn("validation failed for AdminLoginRequest", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}

	tokenData, err := h.service.AuthenticateAdmin(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			slog.Warn("invalid admin credentials", "username", req.Username)
			return echo.ErrUnauthorized.WithInternal(err)
		}
		slog.Error("admin authentication failed", "error", err, "username", req.Username)
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
	return emptyJSON(c, http.StatusOK)
}

func (h *Handler) AdminLogout(c echo.Context) error {
	slog.Debug("AdminLogout handler called", "method", c.Request().Method, "path", c.Path())
	return h.Logout(c) // Reuse the same logout logic
}

// User management handlers
func (h *Handler) GetUsers(c echo.Context) error {
	slog.Debug("GetUsers handler called", "method", c.Request().Method, "path", c.Path())
	users, err := h.service.GetAllUsers()
	if err != nil {
		slog.Error("failed to get users", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) CreateUser(c echo.Context) error {
	slog.Debug("CreateUser handler called", "method", c.Request().Method, "path", c.Path())
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("failed to bind CreateUserRequest", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Warn("validation failed for CreateUserRequest", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}

	user, err := h.service.CreateUser(req)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) GetUser(c echo.Context) error {
	slog.Debug("GetUser handler called", "method", c.Request().Method, "path", c.Path())
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		slog.Warn("invalid user ID", "id", idStr)
		return echo.ErrBadRequest.WithInternal(err)
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		slog.Error("failed to get user", "error", err, "id", id)
		return echo.ErrNotFound.WithInternal(err)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	slog.Debug("UpdateUser handler called", "method", c.Request().Method, "path", c.Path())
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		slog.Warn("invalid user ID", "id", idStr)
		return echo.ErrBadRequest.WithInternal(err)
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("failed to bind UpdateUserRequest", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := h.validator.Struct(req); err != nil {
		slog.Warn("validation failed for UpdateUserRequest", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}

	user, err := h.service.UpdateUser(uint(id), req)
	if err != nil {
		slog.Error("failed to update user", "error", err, "id", id)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	slog.Debug("DeleteUser handler called", "method", c.Request().Method, "path", c.Path())
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		slog.Warn("invalid user ID", "id", idStr)
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		slog.Error("failed to delete user", "error", err, "id", id)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return emptyJSON(c, http.StatusOK)
}

func (h *Handler) GetRefreshTokens(c echo.Context) error {
	slog.Debug("GetRefreshTokens handler called", "method", c.Request().Method, "path", c.Path())
	tokens, err := h.service.GetAllRefreshTokens()
	if err != nil {
		slog.Error("failed to get refresh tokens", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return c.JSON(http.StatusOK, tokens)
}

func (h *Handler) DeleteRefreshToken(c echo.Context) error {
	slog.Debug("DeleteRefreshToken handler called", "method", c.Request().Method, "path", c.Path())
	tokenValue := c.Param("value")
	if tokenValue == "" {
		slog.Warn("empty refresh token value")
		return echo.ErrBadRequest
	}

	if err := h.service.InvalidateRefreshToken(tokenValue); err != nil {
		slog.Error("failed to delete refresh token", "error", err, "value", tokenValue)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	return emptyJSON(c, http.StatusOK)
}
