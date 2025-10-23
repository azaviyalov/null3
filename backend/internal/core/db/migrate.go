package db

import (
	"context"
	"fmt"

	"github.com/azaviyalov/null3/backend/internal/core/auth"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/azaviyalov/null3/backend/internal/mood"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	types := []any{
		&mood.Entry{},
		&auth.RefreshToken{},
	}

	logging.Debug(context.Background(), "attempting database migration")
	if err := db.AutoMigrate(types...); err != nil {
		logging.Error(context.Background(), "database migration failed", "error", err)
		return fmt.Errorf("database migration failed: %w", err)
	}

	logging.Info(context.Background(), "database migration completed successfully")
	return nil
}
