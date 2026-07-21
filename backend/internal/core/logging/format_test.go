package logging_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   logging.Format
		wantName string
	}{
		{name: "text", format: logging.FormatText, wantName: "text"},
		{name: "JSON", format: logging.FormatJSON, wantName: "json"},
		{name: "fancy", format: logging.FormatFancy, wantName: "fancy"},
		{name: "unknown", format: logging.Format(99), wantName: "text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.format.String(); got != tt.wantName {
				t.Errorf("String() = %q, want %q", got, tt.wantName)
			}
		})
	}
}

func TestFormatFromString(t *testing.T) {
	tests := []struct {
		value string
		want  logging.Format
	}{
		{value: "text", want: logging.FormatText},
		{value: "json", want: logging.FormatJSON},
		{value: "fancy", want: logging.FormatFancy},
		{value: "unknown", want: logging.FormatText},
		{value: "", want: logging.FormatText},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := logging.FormatFromString(tt.value); got != tt.want {
				t.Fatalf("FormatFromString(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestFormatNewSLogHandler(t *testing.T) {
	options := &slog.HandlerOptions{Level: slog.LevelWarn}
	tests := []struct {
		name   string
		format logging.Format
		check  func(slog.Handler) bool
	}{
		{name: "text", format: logging.FormatText, check: func(handler slog.Handler) bool {
			_, ok := handler.(*slog.TextHandler)
			return ok
		}},
		{name: "JSON", format: logging.FormatJSON, check: func(handler slog.Handler) bool {
			_, ok := handler.(*slog.JSONHandler)
			return ok
		}},
		{name: "fancy", format: logging.FormatFancy, check: func(handler slog.Handler) bool {
			_, ok := handler.(*logging.FancyHandler)
			return ok
		}},
		{name: "unknown defaults to text", format: logging.Format(99), check: func(handler slog.Handler) bool {
			_, ok := handler.(*slog.TextHandler)
			return ok
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := tt.format.NewSLogHandler(io.Discard, options)

			if !tt.check(handler) {
				t.Fatalf("NewSLogHandler() returned %T", handler)
			}
			if handler.Enabled(t.Context(), slog.LevelInfo) {
				t.Error("handler enables info below configured warn level")
			}
			if !handler.Enabled(t.Context(), slog.LevelWarn) {
				t.Error("handler disables configured warn level")
			}
		})
	}
}
