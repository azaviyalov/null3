package logging

import "log/slog"

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Debug writes a debug message with the logger attached to ctx.
func Debug(ctx any, msg string, args ...any) {
	fromContext(ctx).Debug(msg, args...)
}

// Info writes an informational message with the logger attached to ctx.
func Info(ctx any, msg string, args ...any) {
	fromContext(ctx).Info(msg, args...)
}

// Warn writes a warning with the logger attached to ctx.
func Warn(ctx any, msg string, args ...any) {
	fromContext(ctx).Warn(msg, args...)
}

// Error writes an error message with the logger attached to ctx.
func Error(ctx any, msg string, args ...any) {
	fromContext(ctx).Error(msg, args...)
}

func defaultLogger() Logger {
	return &callerInjector{l: slog.Default()}
}
