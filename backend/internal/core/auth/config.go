package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Production    bool
	JWTSecret     string
	JWTExpiration time.Duration
	SecureCookies bool
}

func GetConfig() (Config, error) {
	config := Config{}

	// Defaults
	config.JWTExpiration = 24 * time.Hour
	config.SecureCookies = false

	if productionParam := os.Getenv("PRODUCTION"); productionParam != "" {
		production, err := strconv.ParseBool(productionParam)
		if err != nil {
			return Config{}, fmt.Errorf("invalid value for PRODUCTION: %v", err)
		}
		config.Production = production
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWTSecret = jwtSecret
	} else {
		if config.Production {
			return Config{}, fmt.Errorf("JWT_SECRET must be set in production mode")
		}
		fmt.Println("Warning: JWT_SECRET is not set, using randomly generated secret")
		jwtSecret, err := generateRandomSecret()
		if err != nil {
			return Config{}, fmt.Errorf("failed to generate JWT_SECRET: %v", err)
		}
		config.JWTSecret = jwtSecret
	}
	if jwtExpirationParam := os.Getenv("JWT_EXPIRATION"); jwtExpirationParam != "" {
		JWTExpiration, err := time.ParseDuration(jwtExpirationParam)
		if err != nil {
			return Config{}, fmt.Errorf("invalid JWT_EXPIRATION: %v", err)
		}
		if JWTExpiration <= 0 {
			return Config{}, fmt.Errorf("JWT_EXPIRATION must be a positive duration")
		}
		config.JWTExpiration = JWTExpiration
	}
	if secureCookiesParam := os.Getenv("SECURE_COOKIES"); secureCookiesParam != "" {
		secureCookies, err := strconv.ParseBool(secureCookiesParam)
		if err != nil {
			return Config{}, fmt.Errorf("invalid SECURE_COOKIES: %v", err)
		}
		config.SecureCookies = secureCookies
	}

	return config, nil
}

func generateRandomSecret() (string, error) {
	const secretLen = 32
	b := make([]byte, secretLen)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random secret: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

type StubUserConfig struct {
	UserID   uint
	Login    string
	Password string
	Email    string
}

func GetStubUserConfig() StubUserConfig {
	fmt.Println("Warning: Using stub user configuration")
	return StubUserConfig{
		UserID:   1,
		Login:    "admin",
		Password: "password",
		Email:    "admin@example.com",
	}
}
