package account

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	PasswordResetTokenExpiration time.Duration
	FrontendURL                  string
}

func GetConfig() (Config, error) {
	config := Config{
		PasswordResetTokenExpiration: time.Hour,
	}

	if resetExpirationParam := os.Getenv("PASSWORD_RESET_TOKEN_EXPIRATION"); resetExpirationParam != "" {
		resetExpiration, err := time.ParseDuration(resetExpirationParam)
		if err != nil {
			return Config{}, fmt.Errorf("invalid PASSWORD_RESET_TOKEN_EXPIRATION: %v", err)
		}
		if resetExpiration <= 0 {
			return Config{}, fmt.Errorf("PASSWORD_RESET_TOKEN_EXPIRATION must be a positive duration")
		}
		config.PasswordResetTokenExpiration = resetExpiration
	}

	return config, nil
}
