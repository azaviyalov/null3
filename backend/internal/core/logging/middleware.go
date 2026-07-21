package logging

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RequestLogger() echo.MiddlewareFunc {
	requestID := middleware.RequestID()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return requestID(func(c echo.Context) error {
			start := time.Now()
			request := c.Request()
			response := c.Response()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			logger := slog.Default().With("request_id", response.Header().Get(echo.HeaderXRequestID))
			path := c.Path()
			if path == "" {
				path = request.URL.Path
			}
			attrs := []any{
				"method", request.Method,
				"path", path,
				"status", response.Status,
				"latency", time.Since(start).String(),
				"ip", c.RealIP(),
			}
			switch {
			case err == nil:
				logger.Info("HTTP request completed", attrs...)
			case response.Status >= 500:
				logger.Error("HTTP request completed with error", append(attrs, "error", err)...)
			default:
				logger.Warn("HTTP request completed with non-fatal error", append(attrs, "error", err)...)
			}
			return err
		})
	}
}
