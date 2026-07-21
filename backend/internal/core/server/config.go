package server

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Address     string
	EnableCORS  bool
	FrontendURL string
}

func GetConfig() (Config, error) {
	config := Config{
		Address:     "localhost:8080",
		FrontendURL: "http://localhost:4200",
	}

	if address := os.Getenv("ADDRESS"); address != "" {
		config.Address = address
	}

	if enableCORSParam := os.Getenv("ENABLE_CORS"); enableCORSParam != "" {
		enable, err := strconv.ParseBool(enableCORSParam)
		if err != nil {
			return config, fmt.Errorf("parse ENABLE_CORS: %w", err)
		}
		config.EnableCORS = enable
	}

	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		config.FrontendURL = frontendURL
	}

	return config, nil
}
