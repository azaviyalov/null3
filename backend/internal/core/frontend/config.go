package frontend

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	EnableFrontendDist bool
	APIURL             string
}

func GetConfig() (Config, error) {
	config := Config{APIURL: "http://localhost:8080/api"}

	if enableFrontendDistParam := os.Getenv("ENABLE_FRONTEND_DIST"); enableFrontendDistParam != "" {
		enableFrontendDist, err := strconv.ParseBool(enableFrontendDistParam)
		if err != nil {
			return config, fmt.Errorf("parse ENABLE_FRONTEND_DIST: %w", err)
		}
		config.EnableFrontendDist = enableFrontendDist
	}

	if config.EnableFrontendDist {
		if apiURL := os.Getenv("API_URL"); apiURL != "" {
			config.APIURL = apiURL
		}
	}

	return config, nil
}
