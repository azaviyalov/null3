package logging

import (
	"context"

	"github.com/labstack/echo/v4"
)

type ctxKey struct{}

var loggerKey = &ctxKey{}

const echoErrorKey = "internal/logging_error"

func fromContext(ctx any) Logger {
	if ctx == nil {
		return defaultLogger()
	}
	switch v := ctx.(type) {
	case context.Context:
		return fromStdContext(v)
	case echo.Context:
		return fromEchoContext(v)
	default:
		panic("unsupported context type")
	}
}

func fromStdContext(ctx context.Context) Logger {
	if ctx == nil {
		return defaultLogger()
	}
	if v := ctx.Value(loggerKey); v != nil {
		if l, ok := v.(Logger); ok && l != nil {
			return l
		}
	}
	return defaultLogger()
}

func fromEchoContext(c echo.Context) Logger {
	if c == nil || c.Request() == nil {
		return defaultLogger()
	}
	ctx := c.Request().Context()
	return fromStdContext(ctx)
}

func withLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func addLogger(c echo.Context, l Logger) {
	if c == nil || c.Request() == nil {
		panic("invalid echo context")
	}
	ctx := c.Request().Context()
	ctx = withLogger(ctx, l)
	c.SetRequest(c.Request().WithContext(ctx))
}
