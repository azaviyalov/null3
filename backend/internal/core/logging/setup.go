package logging

import (
	"log/slog"
	"os"
)

func Setup(config Config) {
	handler := config.Format.NewSLogHandler(os.Stdout, &slog.HandlerOptions{
		Level:       config.Level,
		AddSource:   true,
		ReplaceAttr: ReplaceSourceAttr,
	})

	slog.SetDefault(slog.New(handler))

	slog.Debug("logger setup complete", "level", config.Level, "format", config.Format)
}
