package db

import (
	"context"
	"fmt"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB, types ...any) error {
	logging.Debug(context.Background(), "attempting database migration")
	if err := db.AutoMigrate(types...); err != nil {
		logging.Error(context.Background(), "database migration failed", "error", err)
		return fmt.Errorf("database migration failed: %w", err)
	}

	logging.Info(context.Background(), "database migration completed successfully")
	return nil
}
