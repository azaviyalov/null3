package db

import (
	"github.com/azaviyalov/null3/backend/internal/mood"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	_ = db.AutoMigrate(&mood.Entry{})
}
