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

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}))
	slog.SetDefault(logger)
}
