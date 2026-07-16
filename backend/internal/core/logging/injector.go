package logging

import "log/slog"

type callerInjector struct {
	l *slog.Logger
}

func (a *callerInjector) addCallerFields(args ...any) []any {
	fn := findExternalFuncName()
	if fn == "" {
		return args
	}
	for i := 0; i+1 < len(args); i += 2 {
		if key, ok := args[i].(string); ok && key == "caller" {
			return args
		}
	}

	return append([]any{"caller", fn}, args...)
}

func (a *callerInjector) Debug(msg string, args ...any) {
	a.l.Debug(msg, a.addCallerFields(args...)...)
}

func (a *callerInjector) Info(msg string, args ...any) {
	a.l.Info(msg, a.addCallerFields(args...)...)
}

func (a *callerInjector) Warn(msg string, args ...any) {
	a.l.Warn(msg, a.addCallerFields(args...)...)
}

func (a *callerInjector) Error(msg string, args ...any) {
	a.l.Error(msg, a.addCallerFields(args...)...)
}
