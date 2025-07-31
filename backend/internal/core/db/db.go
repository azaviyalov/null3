package db

import "gorm.io/gorm"

func Setup(config Config) (*gorm.DB, error) {
	db := Connect(config)
	AutoMigrate(db)
	return db, nil
}
