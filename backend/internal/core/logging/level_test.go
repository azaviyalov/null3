package logging_test

import (
	"log/slog"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
)

func TestLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    logging.Level
		wantSlog slog.Level
		wantName string
	}{
		{name: "debug", level: logging.LevelDebug, wantSlog: slog.LevelDebug, wantName: "debug"},
		{name: "info", level: logging.LevelInfo, wantSlog: slog.LevelInfo, wantName: "info"},
		{name: "warn", level: logging.LevelWarn, wantSlog: slog.LevelWarn, wantName: "warn"},
		{name: "error", level: logging.LevelError, wantSlog: slog.LevelError, wantName: "error"},
		{name: "unknown", level: logging.Level(99), wantSlog: slog.LevelInfo, wantName: "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.Level(); got != tt.wantSlog {
				t.Errorf("Level() = %v, want %v", got, tt.wantSlog)
			}
			if got := tt.level.String(); got != tt.wantName {
				t.Errorf("String() = %q, want %q", got, tt.wantName)
			}
		})
	}
}

func TestLevelFromString(t *testing.T) {
	tests := []struct {
		value string
		want  logging.Level
	}{
		{value: "debug", want: logging.LevelDebug},
		{value: "info", want: logging.LevelInfo},
		{value: "warn", want: logging.LevelWarn},
		{value: "error", want: logging.LevelError},
		{value: "unknown", want: logging.LevelInfo},
		{value: "", want: logging.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := logging.LevelFromString(tt.value); got != tt.want {
				t.Fatalf("LevelFromString(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
