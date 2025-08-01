package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	serverErr := make(chan error, 1)
	go func() {
		if err := e.Start(config.Host); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	select {
	case <-quit:
		slog.Info("received shutdown signal, shutting down server gracefully")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
			return err
		}
		slog.Info("server stopped gracefully")
		return nil
	case err := <-serverErr:
		if err != nil {
			slog.Error("server start failed", "error", err)
			return err
		}
		slog.Info("server stopped successfully")
		return nil
	}
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
