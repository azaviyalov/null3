package logging

import (
	"log/slog"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Debug logs a debug-level message using [Logger] extracted from the given context.
// It supports both [context.Context] and [echo.Context].
//
// It returns a default logger if the context is nil.
// It panics if the context is a non-nil value of an unsupported type.
// Be sure to pass a valid context.
func Debug(ctx any, msg string, args ...any) {
	fromContext(ctx).Debug(msg, args...)
}

// Info logs an info-level message using [Logger] extracted from the given context.
// It supports both [context.Context] and [echo.Context].
//
// It returns a default logger if the context is nil.
// It panics if the context is a non-nil value of an unsupported type.
// Be sure to pass a valid context.
func Info(ctx any, msg string, args ...any) {
	fromContext(ctx).Info(msg, args...)
}

// Warn logs a warn-level message using [Logger] extracted from the given context.
// It supports both [context.Context] and [echo.Context].
//
// It returns a default logger if the context is nil.
// It panics if the context is a non-nil value of an unsupported type.
// Be sure to pass a valid context.
func Warn(ctx any, msg string, args ...any) {
	fromContext(ctx).Warn(msg, args...)
}

// Error logs an error-level message using [Logger] extracted from the given context.
// It supports both [context.Context] and [echo.Context].
//
// It returns a default logger if the context is nil.
// It panics if the context is a non-nil value of an unsupported type.
// Be sure to pass a valid context.
func Error(ctx any, msg string, args ...any) {
	fromContext(ctx).Error(msg, args...)
}

func defaultLogger() Logger {
	return &callerInjector{l: slog.Default()}
}
