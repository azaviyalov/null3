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
	config := Config{}

	// Defaults
	config.EnableFrontendDist = true
	config.APIURL = "http://localhost:8080/api"

	if enableFrontendDistParam := os.Getenv("ENABLE_FRONTEND_DIST"); enableFrontendDistParam != "" {
		enableFrontendDist, err := strconv.ParseBool(enableFrontendDistParam)
		if err != nil {
			return config, fmt.Errorf("invalid value for ENABLE_FRONTEND_DIST: %v", err)
		}
		config.EnableFrontendDist = enableFrontendDist
	}

	if enableFrontendDist := os.Getenv("ENABLE_FRONTEND_DIST"); enableFrontendDist != "" {
		if enable, err := strconv.ParseBool(enableFrontendDist); err == nil {
			config.EnableFrontendDist = enable
		}
	}

	if config.EnableFrontendDist {
		if apiURL := os.Getenv("API_URL"); apiURL != "" {
			config.APIURL = apiURL
		}
	} else {
		fmt.Println("Warning: API_URL is not applicable when ENABLE_FRONTEND_DIST is true")
		config.APIURL = ""
	}

	return config, nil
}
