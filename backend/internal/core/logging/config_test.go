package logging_test

import (
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		format string
		want   logging.Config
	}{
		{
			name: "defaults",
			want: logging.Config{Level: logging.LevelInfo, Format: logging.FormatText},
		},
		{
			name:   "environment values",
			level:  "debug",
			format: "json",
			want:   logging.Config{Level: logging.LevelDebug, Format: logging.FormatJSON},
		},
		{
			name:   "unknown values use defaults",
			level:  "verbose",
			format: "yaml",
			want:   logging.Config{Level: logging.LevelInfo, Format: logging.FormatText},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("LOG_LEVEL", tt.level)
			t.Setenv("LOG_FORMAT", tt.format)

			got := logging.GetConfig()

			if got != tt.want {
				t.Fatalf("GetConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
