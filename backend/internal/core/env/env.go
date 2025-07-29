package env

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Setup() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	// Set config defaults if not set
	if os.Getenv("ENV") == "" {
		os.Setenv("ENV", "production")
	}
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "file:null3.db?_fk=1")
	}
	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "info")
	}
	if os.Getenv("LOG_FORMAT") == "" {
		os.Setenv("LOG_FORMAT", "json")
	}
	if os.Getenv("ENABLE_FRONTEND_DIST") == "" {
		os.Setenv("ENABLE_FRONTEND_DIST", "false")
	}
	if os.Getenv("ENABLE_CORS") == "" {
		os.Setenv("ENABLE_CORS", "false")
	}
	if os.Getenv("FRONTEND_URL") == "" {
		os.Setenv("FRONTEND_URL", "http://localhost:4200")
	}
	if os.Getenv("API_URL") == "" {
		os.Setenv("API_URL", fmt.Sprintf("http://localhost:%s/api", os.Getenv("PORT")))
	}
	if os.Getenv("JWT_SECRET") == "" {
		if os.Getenv("ENV") == "production" {
			fmt.Println("Error: JWT_SECRET must be set in production environments.")
			os.Exit(1)
		}

		fmt.Println("Warning: JWT_SECRET not set, generating a random secret for development")
		os.Setenv("JWT_SECRET", generateRandomSecret())
	}
	if os.Getenv("JWT_EXPIRATION") == "" {
		os.Setenv("JWT_EXPIRATION", "24h")
	}

	// These are placeholders until user management is implemented
	if os.Getenv("USER_ID") == "" {
		os.Setenv("USER_ID", "1")
	}
	if os.Getenv("LOGIN") == "" {
		os.Setenv("LOGIN", "admin")
	}
	if os.Getenv("PASSWORD") == "" {
		os.Setenv("PASSWORD", "password")
	}
	if os.Getenv("EMAIL") == "" {
		os.Setenv("EMAIL", "admin@example.com")
	}
}

func generateRandomSecret() string {
	const secretLen = 32
	b := make([]byte, secretLen)
	_, err := rand.Read(b)
	if err != nil {
		// fallback: return a static string if random fails
		return "fallback_secret"
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
