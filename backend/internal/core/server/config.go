package server

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host        string
	EnableCORS  bool
	FrontendURL string
}

func GetConfig() (Config, error) {
	config := Config{}

	// Defaults
	config.Host = "localhost:8080"
	config.EnableCORS = false
	config.FrontendURL = "http://localhost:4200" // Default frontend URL when CORS is enabled

	if host := os.Getenv("HOST"); host != "" {
		config.Host = host
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
	} else {
		config.FrontendURL = "unused" // Placeholder when CORS is not enabled
	}

	return config, nil
}
