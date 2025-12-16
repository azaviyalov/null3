package logging

import (
	"context"

	"github.com/labstack/echo/v4"
)

type ctxKey struct{}

var loggerKey = &ctxKey{}

const echoErrorKey = "internal/logging_error"

// fromContext extracts [Logger] from various context types.
// It supports both [context.Context] and [echo.Context].
// It returns a default logger if the context is nil or no logger is found in the context.
// It panics if the context is a non-nil value of an unsupported type.
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

// fromStdContext extracts [Logger] from [context.Context].
// It returns a default logger if the context is nil or no logger is found in the context.
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

// fromEchoContext extracts [Logger] from [echo.Context].
// It returns a default logger if the echo context, its request, or its request context is nil.
func fromEchoContext(c echo.Context) Logger {
	if c == nil || c.Request() == nil {
		return defaultLogger()
	}
	ctx := c.Request().Context()
	return fromStdContext(ctx)
}

// withLogger returns a new [context.Context] with the given [Logger] attached.
func withLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// addLogger adds the given [Logger] to the [echo.Context]'s underlying [context.Context].
// It panics if the echo context or its request is nil.
func addLogger(c echo.Context, l Logger) {
	if c == nil || c.Request() == nil {
		panic("invalid echo context")
	}
	ctx := c.Request().Context()
	ctx = withLogger(ctx, l)
	c.SetRequest(c.Request().WithContext(ctx))
}
