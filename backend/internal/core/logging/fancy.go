package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

type FancyHandler struct {
	level        slog.Leveler
	addSource    bool
	replaceAttrs func(groups []string, a slog.Attr) slog.Attr
	attrs        []slog.Attr
}

func NewFancyHandler(options *slog.HandlerOptions) *FancyHandler {
	return &FancyHandler{
		level:        options.Level,
		addSource:    options.AddSource,
		replaceAttrs: options.ReplaceAttr,
	}
}

func (h *FancyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *FancyHandler) Handle(_ context.Context, r slog.Record) error {
	timeStr := fmt.Sprintf("\033[35m%s\033[0m", time.Now().Format("15:04:05.000"))
	lvl := coloredLevel(r.Level)
	msg := fmt.Sprintf("\033[1m%s\033[0m", r.Message)
	caller := h.caller()
	attrs := h.formatAttrs(r)
	status := fmt.Sprintf("\033[1m%s %s\033[0m", lvl, caller)
	fmt.Fprintf(os.Stdout, "%s %s %s %s%s\n", timeStr, status, msg, "", attrs)
	return nil
}

func (h *FancyHandler) caller() string {
	caller := ""
	if h.addSource {
		if _, file, line, ok := runtime.Caller(4); ok {
			shortFile := trimProjectRoot(file)
			caller = fmt.Sprintf("%s:%d", shortFile, line)
		}
	}
	return caller
}

func (h *FancyHandler) formatAttrs(r slog.Record) string {
	var attrLines []string
	for _, a := range h.attrs {
		attrLines = append(attrLines, h.formatAttr(a))
	}
	r.Attrs(func(a slog.Attr) bool {
		attrLines = append(attrLines, h.formatAttr(a))
		return true
	})
	if len(attrLines) > 0 {
		return "\n" + strings.Join(attrLines, "\n")
	}
	return ""
}

func (h *FancyHandler) formatAttr(a slog.Attr) string {
	if h.replaceAttrs != nil {
		a = h.replaceAttrs(nil, a)
	}
	val := a.Value
	if lv, ok := val.Any().(slog.LogValuer); ok {
		val = lv.LogValue()
	}
	key := fmt.Sprintf("\033[1;34m%s\033[0m", a.Key)
	valStr := fmt.Sprintf("%v", val)
	if s, ok := val.Any().(string); ok {
		valStr = fmt.Sprintf("\"%s\"", s)
	}
	return fmt.Sprintf("  - %s: %s", key, valStr)
}

func (h *FancyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &FancyHandler{
		level:        h.level,
		addSource:    h.addSource,
		replaceAttrs: h.replaceAttrs,
		attrs:        newAttrs,
	}
}

func (h *FancyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	// Prepend group name to all current attrs
	newAttrs := make([]slog.Attr, len(h.attrs))
	for i, a := range h.attrs {
		a.Key = name + "." + a.Key
		newAttrs[i] = a
	}
	return &FancyHandler{
		level:        h.level,
		addSource:    h.addSource,
		replaceAttrs: h.replaceAttrs,
		attrs:        newAttrs,
	}
}

func coloredLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "\033[36m[DEBU]\033[0m"
	case slog.LevelInfo:
		return "\033[32m[INFO]\033[0m"
	case slog.LevelWarn:
		return "\033[33m[WARN]\033[0m"
	case slog.LevelError:
		return "\033[31m[ERRO]\033[0m"
	default:
		return fmt.Sprintf("%5s", level.String())
	}
}

func trimProjectRoot(file string) string {
	marker := "null3/backend/"
	if idx := strings.Index(file, marker); idx != -1 {
		return file[idx:]
	}
	return file
}
