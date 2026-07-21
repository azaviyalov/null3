package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
)

type FancyHandler struct {
	writer       io.Writer
	level        slog.Leveler
	addSource    bool
	replaceAttrs func(groups []string, a slog.Attr) slog.Attr
	attrs        []fancyAttr
	groups       []string
}

type fancyAttr struct {
	groups []string
	attr   slog.Attr
}

func NewFancyHandler(writer io.Writer, options *slog.HandlerOptions) *FancyHandler {
	return &FancyHandler{
		writer:       writer,
		level:        options.Level,
		addSource:    options.AddSource,
		replaceAttrs: options.ReplaceAttr,
	}
}

func (h *FancyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *FancyHandler) Handle(_ context.Context, r slog.Record) error {
	timeStr := fmt.Sprintf("\033[35m%s\033[0m", r.Time.Format("15:04:05.000"))
	lvl := coloredLevel(r.Level)
	msg := fmt.Sprintf("\033[1m%s\033[0m", r.Message)
	attrs := h.formatRecordAttrs(r)
	status := fmt.Sprintf("\033[1m%s\033[0m", lvl)
	_, err := fmt.Fprintf(h.writer, "%s %s %s%s\n", timeStr, status, msg, attrs)
	return err
}

func (h *FancyHandler) formatRecordAttrs(r slog.Record) string {
	var attrLines []string
	for _, stored := range h.attrs {
		attrLines = append(attrLines, h.formatAttr(stored.groups, stored.attr))
	}
	if h.addSource {
		if fileLine := sourceFromPC(r.PC); fileLine != "" {
			attrLines = append(attrLines, h.formatAttr(nil, slog.String("source", fileLine)))
		}
	}
	r.Attrs(func(a slog.Attr) bool {
		attrLines = append(attrLines, h.formatAttr(h.groups, a))
		return true
	})
	if len(attrLines) > 0 {
		return "\n" + strings.Join(attrLines, "\n")
	}
	return ""
}

func (h *FancyHandler) formatAttr(groups []string, a slog.Attr) string {
	if h.replaceAttrs != nil {
		a = h.replaceAttrs(groups, a)
	}
	val := a.Value.Resolve()
	key := a.Key
	if len(groups) > 0 {
		key = strings.Join(groups, ".") + "." + key
	}
	key = fmt.Sprintf("\033[1;34m%s\033[0m", key)
	valStr := fmt.Sprintf("%v", val)
	if s, ok := val.Any().(string); ok {
		valStr = strconv.Quote(s)
	}
	return fmt.Sprintf("  - %s: %s", key, valStr)
}

func (h *FancyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]fancyAttr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	for _, attr := range attrs {
		newAttrs = append(newAttrs, fancyAttr{
			groups: append([]string(nil), h.groups...),
			attr:   attr,
		})
	}
	return &FancyHandler{
		writer:       h.writer,
		level:        h.level,
		addSource:    h.addSource,
		replaceAttrs: h.replaceAttrs,
		attrs:        newAttrs,
		groups:       h.groups,
	}
}

func (h *FancyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &FancyHandler{
		writer:       h.writer,
		level:        h.level,
		addSource:    h.addSource,
		replaceAttrs: h.replaceAttrs,
		attrs:        h.attrs,
		groups:       append(append([]string(nil), h.groups...), name),
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
