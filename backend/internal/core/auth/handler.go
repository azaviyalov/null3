package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID uint `json:"id"`
}

func LoginHandler(c echo.Context) error {
	slog.Debug("LoginHandler called", "method", c.Request().Method, "path", c.Path())
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("failed to bind LoginRequest", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := c.Validate(req); err != nil {
		slog.Warn("validation failed for LoginRequest", "error", err)
		return echo.ErrBadRequest.WithInternal(err)
	}
	// For simplicity, we use a stub function to check credentials
	if err := checkLoginCredentialsStub(req); err != nil {
		slog.Warn("check login credentials failed", "login", req.Login, "error", err)
		if errors.Is(err, ErrInvalidCredentials) {
			return echo.ErrUnauthorized.WithInternal(err)
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	// For simplicity, we use a stub function to generate a token
	token, err := generateTokenStub()
	if err != nil {
		slog.Error("token generation failed", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}

	slog.Info("User logged in", "login", req.Login)
	// For simplicity, we use a stub user ID from environment variable
	userIDStr := os.Getenv("USER_ID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		slog.Error("failed to parse USER_ID", "error", err)
		return echo.ErrInternalServerError.WithInternal(err)
	}
	response := UserResponse{
		ID: uint(userID),
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt_token"
	cookie.Value = token
	cookie.HttpOnly = true
	if os.Getenv("ENV") == "production" {
		cookie.Secure = true
	}
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, response)
}

func checkLoginCredentialsStub(req LoginRequest) error {
	login := os.Getenv("LOGIN")
	password := os.Getenv("PASSWORD")
	// Allow case-insensitive login, trim spaces
	if strings.TrimSpace(strings.ToLower(req.Login)) != strings.TrimSpace(strings.ToLower(login)) ||
		req.Password != password {
		return ErrInvalidCredentials
	}
	return nil
}

func generateTokenStub() (string, error) {
	userID := os.Getenv("USER_ID")
	token, err := GenerateJWT(userID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenGenerationFailed, err)
	}
	return token, nil
}

func MeHandler(c echo.Context) error {
	user, _ := GetUser(c)
	meResponse := UserResponse{
		ID: user.ID,
	}
	return c.JSON(http.StatusOK, meResponse)
}

func LogoutHandler(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt_token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = -1
	cookie.SameSite = http.SameSiteLaxMode
	if os.Getenv("ENV") == "production" {
		cookie.Secure = true
	}
	c.SetCookie(cookie)
	return c.NoContent(http.StatusNoContent)
}
