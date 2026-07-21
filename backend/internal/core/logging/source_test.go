package logging_test

import (
	"log/slog"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
)

func TestReplaceSourceAttr(t *testing.T) {
	t.Run("formats source", func(t *testing.T) {
		attr := slog.Any("source", &slog.Source{
			File: "/workspace/null3/backend/internal/app/app.go",
			Line: 42,
		})

		got := logging.ReplaceSourceAttr(nil, attr)

		if got.Key != "source" {
			t.Errorf("key = %q, want %q", got.Key, "source")
		}
		if got.Value.Kind() != slog.KindString {
			t.Fatalf("value kind = %v, want string", got.Value.Kind())
		}
		if got.Value.String() != "backend/internal/app/app.go:42" {
			t.Fatalf("value = %q, want %q", got.Value.String(), "backend/internal/app/app.go:42")
		}
	})

	tests := []struct {
		name string
		attr slog.Attr
	}{
		{name: "unrelated attribute", attr: slog.String("component", "server")},
		{name: "source with wrong value type", attr: slog.String("source", "already formatted")},
		{name: "nil source", attr: slog.Any("source", (*slog.Source)(nil))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logging.ReplaceSourceAttr(nil, tt.attr); !got.Equal(tt.attr) {
				t.Fatalf("ReplaceSourceAttr() = %v, want unchanged %v", got, tt.attr)
			}
		})
	}
}
