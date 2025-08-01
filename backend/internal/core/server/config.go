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
	config := Config{}

	// Defaults
	config.Address = "localhost:8080"
	config.EnableCORS = false
	config.FrontendURL = "http://localhost:4200" // Default frontend URL (used when CORS is enabled)

	if address := os.Getenv("ADDRESS"); address != "" {
		config.Address = address
	}

	if enableCORSParam := os.Getenv("ENABLE_CORS"); enableCORSParam != "" {
		enable, err := strconv.ParseBool(enableCORSParam)
		if err != nil {
			return config, fmt.Errorf("invalid value for ENABLE_CORS: %v", err)
		}
		config.EnableCORS = enable
	}

	if config.EnableCORS {
		if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
			config.FrontendURL = frontendURL
		}
	}

	return config, nil
}
