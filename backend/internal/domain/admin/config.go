package admin

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	FrontendURL string
	Password    string
}

func GetConfig() (Config, error) {
	password := os.Getenv("ADMIN_PASSWORD")
	if strings.TrimSpace(password) == "" {
		return Config{}, fmt.Errorf("ADMIN_PASSWORD must be set and non-empty")
	}
	return Config{Password: password}, nil
}
