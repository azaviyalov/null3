package env

import (
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
	if os.Getenv("API_URL") == "" {
		os.Setenv("API_URL", fmt.Sprintf("http://localhost:%s/api", os.Getenv("PORT")))
	}
	if os.Getenv("JWT_SECRET") == "" {
		if os.Getenv("ENV") == "production" {
			fmt.Println("Error: JWT_SECRET must be set in production environments.")
			os.Exit(1)
		}
		// Default value for non-production environments
		os.Setenv("JWT_SECRET", "example_secret")
	}
	if os.Getenv("JWT_EXPIRATION") == "" {
		os.Setenv("JWT_EXPIRATION", "24h")
	}

	// These are placeholders until user management is implemented
	if os.Getenv("USER_ID") == "" {
		os.Setenv("USER_ID", "1")
	}
	if os.Getenv("LOGIN") == "" {
		os.Setenv("LOGIN", "login")
	}
	if os.Getenv("PASSWORD") == "" {
		os.Setenv("PASSWORD", "password")
	}
	if os.Getenv("EMAIL") == "" {
		os.Setenv("EMAIL", "user@example.com")
	}
}
