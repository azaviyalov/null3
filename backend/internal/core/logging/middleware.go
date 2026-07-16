package logging

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

const echoErrorKey = "internal/logging_error"

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			reqID := c.Request().Header.Get("X-Request-Id")
			if reqID == "" {
				reqID = generateRequestID()
			}
			c.Response().Header().Set("X-Request-Id", reqID)
			reqLogger := slog.Default().With("request_id", reqID)
			c.Response().After(func() {
				req := c.Request()
				res := c.Response()
				status := res.Status
				attrs := []any{
					"method", req.Method,
					"path", req.URL.Path,
					"status", status,
					"latency", time.Since(start).String(),
					"ip", c.RealIP(),
				}
				if v := c.Get(echoErrorKey); v != nil {
					if status >= 500 {
						reqLogger.Error("HTTP request completed with error", append(attrs, "error", v)...)
					} else {
						reqLogger.Warn("HTTP request completed with non-fatal error", append(attrs, "error", v)...)
					}
					return
				}
				reqLogger.Info("HTTP request completed", attrs...)
			})
			err := next(c)
			if err != nil {
				c.Set(echoErrorKey, err)
			}
			return err
		}
	}
}

func generateRequestID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
