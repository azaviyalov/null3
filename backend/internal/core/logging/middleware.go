package logging

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			reqID := c.Request().Header.Get("X-Request-Id")
			if reqID == "" {
				reqID = generateRequestID()
			}
			c.Response().Header().Set("X-Request-Id", reqID)
			baseReqLogger := slog.Default().With("request_id", reqID)
			reqLogger := &callerInjector{l: baseReqLogger}
			withEchoLogger(c, reqLogger)
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
						ErrorEcho(c, "HTTP request completed with error", append(attrs, "error", v)...)
					} else {
						WarnEcho(c, "HTTP request completed with non-fatal error", append(attrs, "error", v)...)
					}
					return
				}
				InfoEcho(c, "HTTP request completed", attrs...)
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
		// fallback to timestamp if random fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
