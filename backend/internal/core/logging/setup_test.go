package logging_test

import (
	"log/slog"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
)

func TestSetupInstallsConfiguredDefaultLogger(t *testing.T) {
	previous := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(previous)
	})

	logging.Setup(logging.Config{Level: logging.LevelError, Format: logging.FormatFancy})

	handler := slog.Default().Handler()
	if _, ok := handler.(*logging.FancyHandler); !ok {
		t.Fatalf("default handler = %T, want *logging.FancyHandler", handler)
	}
	if handler.Enabled(t.Context(), slog.LevelWarn) {
		t.Error("default handler enables warn below configured error level")
	}
	if !handler.Enabled(t.Context(), slog.LevelError) {
		t.Error("default handler disables configured error level")
	}
}
