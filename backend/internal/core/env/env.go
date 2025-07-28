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
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "file:null3.db?_fk=1")
	}
	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "info")
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
}
