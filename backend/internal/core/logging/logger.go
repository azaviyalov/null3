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

	logFormat := os.Getenv("LOG_FORMAT")
	var loggerHandler slog.Handler
	switch logFormat {
	case "json":
		loggerHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	case "text":
		loggerHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	case "fancy":
		loggerHandler = NewFancyHandler(&slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	default:
		fmt.Printf("invalid LOG_FORMAT: %s, must be one of json, text, fancy\n", logFormat)
		os.Exit(1)
	}
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)
	slog.Debug("logger setup complete", "level", logLevel, "format", logFormat)
}
