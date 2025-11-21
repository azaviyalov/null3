package logging

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Debug logs a debug-level message using [Logger] extracted from the given context.
// It supports both [context.Context] and [echo.Context].
// It panics if the context type is unsupported or if the context is nil.
func Debug(ctx any, msg string, args ...any) {
	fromContext(ctx).Debug(msg, args...)
}

// Info logs an info-level message using [Logger] extracted from the given context.
// It supports both [context.Context] and [echo.Context].
// It panics if the context type is unsupported or if the context is nil.
func Info(ctx any, msg string, args ...any) {
	fromContext(ctx).Info(msg, args...)
}

// Warn logs a warn-level message using [Logger] extracted from the given context.
// It supports both [context.Context] and [echo.Context].
// It panics if the context type is unsupported or if the context is nil.
func Warn(ctx any, msg string, args ...any) {
	fromContext(ctx).Warn(msg, args...)
}

// Error logs an error-level message using [Logger] extracted from the given context.
// It supports both [context.Context] and [echo.Context].
// It panics if the context type is unsupported or if the context is nil.
func Error(ctx any, msg string, args ...any) {
	fromContext(ctx).Error(msg, args...)
}
