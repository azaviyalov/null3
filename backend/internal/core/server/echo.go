package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	slog.Info("starting HTTP server", "address", config.Address)

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- e.Start(config.Address)
	}()

	select {
	case <-shutdownSignal.Done():
		slog.Info("received shutdown signal, shutting down server gracefully")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shut down HTTP server: %w", err)
		}
		slog.Info("server stopped gracefully")
		return nil
	case err := <-serverErr:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("start HTTP server: %w", err)
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
