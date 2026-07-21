package logging_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
)

func TestFancyHandlerEnabled(t *testing.T) {
	handler := logging.NewFancyHandler(&bytes.Buffer{}, &slog.HandlerOptions{Level: slog.LevelWarn})

	if handler.Enabled(t.Context(), slog.LevelInfo) {
		t.Error("Enabled(info) = true, want false")
	}
	if !handler.Enabled(t.Context(), slog.LevelWarn) {
		t.Error("Enabled(warn) = false, want true")
	}
}

func TestFancyHandlerFormatsAttributes(t *testing.T) {
	output := &bytes.Buffer{}
	handler := logging.NewFancyHandler(output, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == "renamed" {
				attr.Key = "replacement"
			}
			return attr
		},
	})
	derived := handler.WithAttrs([]slog.Attr{
		slog.String("component", "server\nprimary"),
		slog.Int("attempt", 2),
	}).WithGroup("request").WithAttrs([]slog.Attr{
		slog.String("phase", "done"),
	})
	record := slog.NewRecord(time.Time{}, slog.LevelInfo, "handled", 0)
	record.AddAttrs(slog.String("renamed", "yes"))

	if err := derived.Handle(t.Context(), record); err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	for _, want := range []string{
		"component",
		`"server\nprimary"`,
		"attempt",
		"request.phase",
		"request.replacement",
		"\x22yes\x22",
	} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("output %q does not contain %q", output.String(), want)
		}
	}
	if strings.Contains(output.String(), "request.component") {
		t.Errorf("output %q moved pre-group attribute into group", output.String())
	}
}

func TestFancyHandlerWritesFormattedRecord(t *testing.T) {
	recordTime := time.Date(2026, time.July, 18, 12, 34, 56, 789_000_000, time.UTC)
	tests := []struct {
		name  string
		level slog.Level
		want  string
	}{
		{name: "debug", level: slog.LevelDebug, want: "\x1b[36m[DEBU]\x1b[0m"},
		{name: "info", level: slog.LevelInfo, want: "\x1b[32m[INFO]\x1b[0m"},
		{name: "warn", level: slog.LevelWarn, want: "\x1b[33m[WARN]\x1b[0m"},
		{name: "error", level: slog.LevelError, want: "\x1b[31m[ERRO]\x1b[0m"},
		{name: "custom", level: slog.Level(2), want: "INFO+2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			handler := logging.NewFancyHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug})
			record := slog.NewRecord(recordTime, tt.level, "request completed", 0)
			record.AddAttrs(slog.Int("status", 204))

			if err := handler.Handle(t.Context(), record); err != nil {
				t.Fatalf("Handle() error = %v", err)
			}

			for _, want := range []string{
				"12:34:56.789",
				tt.want,
				"\x1b[1mrequest completed\x1b[0m",
				"status",
				"204",
			} {
				if !strings.Contains(output.String(), want) {
					t.Errorf("output %q does not contain %q", output.String(), want)
				}
			}
		})
	}
}
