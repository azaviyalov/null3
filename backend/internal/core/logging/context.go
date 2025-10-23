package logging

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type ctxKey struct{}

var loggerKey = &ctxKey{}

const echoErrorKey = "internal/logging_error"

func withLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func fromContext(ctx context.Context) Logger {
	if ctx == nil {
		return &callerInjector{l: slog.Default()}
	}
	if v := ctx.Value(loggerKey); v != nil {
		if l, ok := v.(Logger); ok && l != nil {
			return l
		}
	}
	return &callerInjector{l: slog.Default()}
}

func fromEcho(c echo.Context) Logger {
	if c == nil || c.Request() == nil {
		return fromContext(context.TODO())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.TODO()
	}
	return fromContext(ctx)
}

func withEchoLogger(c echo.Context, l Logger) {
	if c == nil || c.Request() == nil {
		return
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.TODO()
	}
	ctx = withLogger(ctx, l)
	c.SetRequest(c.Request().WithContext(ctx))
}
