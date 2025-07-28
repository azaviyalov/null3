package db

import (
	"log/slog"
	"os"

	"github.com/azaviyalov/null3/backend/internal/mood"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	types := []any{
		&mood.Entry{},
	}

	slog.Debug("attempting database migration")
	if err := db.AutoMigrate(types...); err != nil {
		slog.Error("database migration failed", "error", err)
		os.Exit(1)
	}

	slog.Info("database migration completed successfully")
}
