package logging

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func Debug(ctx any, msg string, args ...any) {
	fromContext(ctx).Debug(msg, args...)
}

func Info(ctx any, msg string, args ...any) {
	fromContext(ctx).Info(msg, args...)
}

func Warn(ctx any, msg string, args ...any) {
	fromContext(ctx).Warn(msg, args...)
}

func Error(ctx any, msg string, args ...any) {
	fromContext(ctx).Error(msg, args...)
}
