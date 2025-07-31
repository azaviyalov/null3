package logging

import (
	"log/slog"
)

func Setup(config Config) {
	handler := config.Format.NewSLogHandler(&slog.HandlerOptions{
		Level:     config.Level,
		AddSource: true,
	})

	slog.SetDefault(slog.New(handler))

	slog.Debug("logger setup complete", "level", config.Level, "format", config.Format)
}
