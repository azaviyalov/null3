package logging

import (
	"fmt"
	"log/slog"
	"os"
)

func Setup() {
	logLevel := os.Getenv("LOG_LEVEL")
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		fmt.Printf("invalid LOG_LEVEL: %s, must be one of debug, info, warn, error\n", logLevel)
		os.Exit(1)
	}

	if os.Getenv("LOG_FANCY") == "true" {
		logger := slog.New(&FancyHandler{level: level, AddSource: true})
		slog.SetDefault(logger)
		slog.Debug("fancy logger setup complete", "level", logLevel)
	} else {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		}))
		slog.SetDefault(logger)
		slog.Debug("logger setup complete", "level", logLevel)
	}
}
