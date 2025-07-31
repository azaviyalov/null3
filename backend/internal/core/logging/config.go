package logging

import "os"

type Config struct {
	Level  Level
	Format Format
}

func GetConfig() Config {
	config := Config{}

	// Defaults
	config.Level = LevelInfo
	config.Format = FormatText

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Level = LevelFromString(level)
	}

	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.Format = FormatFromString(format)
	}

	return config
}
