package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(config Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DatabaseURL), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	return db, nil
}
