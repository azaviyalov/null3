package logging

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	// сontextKeyError is the key used to store error in context
	сontextKeyError = "logging_middleware_error"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			// Register an after hook to log after response is committed
			c.Response().After(func() {
				req := c.Request()
				res := c.Response()
				status := res.Status
				slog.Info("http_request",
					"method", req.Method,
					"path", req.URL.Path,
					"status", status,
					"latency", time.Since(start).String(),
					"ip", c.RealIP(),
					"error", c.Get(сontextKeyError),
				)
			})
			err := next(c)
			// Save error in context for after hook
			if err != nil {
				c.Set(сontextKeyError, err)
			}
			return err
		}
	}
}
