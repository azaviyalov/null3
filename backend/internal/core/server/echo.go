package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const shutdownTimeout = 10 * time.Second

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
	logging.Info(context.Background(), "starting HTTP server", "address", config.Address)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		defer close(serverErr)
		if err := e.Start(config.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-quit:
		logging.Info(context.Background(), "received shutdown signal, shutting down server gracefully")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			logging.Error(context.Background(), "graceful shutdown failed", "error", err)
			return err
		}
		logging.Info(context.Background(), "server stopped gracefully")
		return nil
	case err := <-serverErr:
		logging.Error(context.Background(), "server start failed", "error", err)
		return err
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
