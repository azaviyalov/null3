package db

import "gorm.io/gorm"

func Setup(config Config) (*gorm.DB, error) {
	db, err := Connect(config)
	if err != nil {
		return nil, err
	}
	if err := AutoMigrate(db); err != nil {
		return nil, err
	}
	return db, nil
}
