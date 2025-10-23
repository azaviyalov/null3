package logging

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func Debug(ctx context.Context, msg string, args ...any) {
	fromContext(ctx).Debug(msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	fromContext(ctx).Info(msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	fromContext(ctx).Warn(msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	fromContext(ctx).Error(msg, args...)
}

func DebugEcho(c echo.Context, msg string, args ...any) {
	fromEcho(c).Debug(msg, args...)
}

func InfoEcho(c echo.Context, msg string, args ...any) {
	fromEcho(c).Info(msg, args...)
}

func WarnEcho(c echo.Context, msg string, args ...any) {
	fromEcho(c).Warn(msg, args...)
}

func ErrorEcho(c echo.Context, msg string, args ...any) {
	fromEcho(c).Error(msg, args...)
}
