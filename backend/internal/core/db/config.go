package db

import "os"

type Config struct {
	DatabaseURL string
}

func GetConfig() Config {
	config := Config{}

	// Defaults
	config.DatabaseURL = "file:null3.db?_fk=1"

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.DatabaseURL = dbURL
	}
	return config
}
