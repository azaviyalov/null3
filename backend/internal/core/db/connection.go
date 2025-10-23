package db

import (
	"context"
	"fmt"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(config Config) (*gorm.DB, error) {
	logging.Debug(context.Background(), "attempting database connection")
	db, err := gorm.Open(sqlite.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		logging.Error(context.Background(), "failed to open database connection", "error", err)
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	logging.Info(context.Background(), "database connection established")
	return db, nil
}
