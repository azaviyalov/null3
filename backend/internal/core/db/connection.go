package db

import (
	"log/slog"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(config Config) *gorm.DB {
	slog.Debug("attempting database connection", "database_url", config.DatabaseURL)
	db, err := gorm.Open(sqlite.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	slog.Info("database connection established")
	return db
}
