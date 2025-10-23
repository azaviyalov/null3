package logging

import "log/slog"

type callerInjector struct {
	l *slog.Logger
}

func (a *callerInjector) addCallerFields(args ...any) []any {
	// compute function name
	fn := findExternalFuncName()
	if fn == "" {
		return args
	}
	// avoid duplicating keys if caller already present in args
	hasKey := func(key string) bool {
		for i := 0; i+1 < len(args); i += 2 {
			if k, ok := args[i].(string); ok && k == key {
				return true
			}
		}
		return false
	}
	newArgs := make([]any, 0, len(args)+2)
	// Only inject caller here. Handlers produce the canonical "source" attr.
	if !hasKey("caller") && fn != "" {
		newArgs = append(newArgs, "caller", fn)
	}
	newArgs = append(newArgs, args...)
	return newArgs
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
