package session

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	JWTSecret              string
	JWTExpiration          time.Duration
	SecureCookies          bool
	RefreshTokenExpiration time.Duration
}

func GetConfig() (Config, error) {
	config := Config{
		JWTExpiration:          24 * time.Hour,
		RefreshTokenExpiration: 7 * 24 * time.Hour,
	}

	config.JWTSecret = os.Getenv("JWT_SECRET")
	if strings.TrimSpace(config.JWTSecret) == "" {
		return Config{}, fmt.Errorf("JWT_SECRET must be set and non-empty")
	}

	if jwtExpirationParam := os.Getenv("JWT_EXPIRATION"); jwtExpirationParam != "" {
		jwtExpiration, err := time.ParseDuration(jwtExpirationParam)
		if err != nil {
			return Config{}, fmt.Errorf("parse JWT_EXPIRATION: %w", err)
		}
		if jwtExpiration <= 0 {
			return Config{}, fmt.Errorf("JWT_EXPIRATION must be a positive duration")
		}
		config.JWTExpiration = jwtExpiration
	}

	if refreshExpirationParam := os.Getenv("REFRESH_TOKEN_EXPIRATION"); refreshExpirationParam != "" {
		refreshExpiration, err := time.ParseDuration(refreshExpirationParam)
		if err != nil {
			return Config{}, fmt.Errorf("parse REFRESH_TOKEN_EXPIRATION: %w", err)
		}
		if refreshExpiration <= 0 {
			return Config{}, fmt.Errorf("REFRESH_TOKEN_EXPIRATION must be a positive duration")
		}
		config.RefreshTokenExpiration = refreshExpiration
	}

	if secureCookiesParam := os.Getenv("SECURE_COOKIES"); secureCookiesParam != "" {
		secureCookies, err := strconv.ParseBool(secureCookiesParam)
		if err != nil {
			return Config{}, fmt.Errorf("parse SECURE_COOKIES: %w", err)
		}
		config.SecureCookies = secureCookies
	}

	return config, nil
}
