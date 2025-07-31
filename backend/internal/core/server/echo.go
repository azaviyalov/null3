package server

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEchoServer(config Config) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	if config.EnableCORS {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{config.FrontendURL},
			AllowCredentials: true,
		}))
	}

	e.Use(logging.RequestLogger())
	e.Use(middleware.Recover())

	e.Validator = newCustomValidator()

	return e
}

func StartServer(e *echo.Echo, config Config) error {
	slog.Info("starting HTTP server", "host", config.Host)
	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server start failed", "error", err)
			return err
		}
	}
	slog.Info("server stopped successfully")
	return nil
}

type customValidator struct {
	v *validator.Validate
}

func newCustomValidator() *customValidator {
	return &customValidator{v: validator.New(validator.WithRequiredStructEnabled())}
}

func (cv *customValidator) Validate(i any) error {
	return cv.v.Struct(i)
}
