package testutil

import (
	"io"
	"log/slog"
	"testing"
)

func DiscardLogs(t testing.TB) {
	t.Helper()
	previousLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	t.Cleanup(func() {
		slog.SetDefault(previousLogger)
	})
}
