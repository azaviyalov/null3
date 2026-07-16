package db

import (
	"fmt"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB, types ...any) error {
	if err := db.AutoMigrate(types...); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}
	return nil
}
