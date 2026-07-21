package app_test

import (
	"testing"

	"github.com/azaviyalov/null3/backend/internal/app"
)

func TestGetConfigPropagatesFrontendURL(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-signing-secret")
	t.Setenv("JWT_EXPIRATION", "")
	t.Setenv("REFRESH_TOKEN_EXPIRATION", "")
	t.Setenv("SECURE_COOKIES", "")
	t.Setenv("PASSWORD_RESET_TOKEN_EXPIRATION", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("ENABLE_FRONTEND_DIST", "")
	t.Setenv("API_URL", "")
	t.Setenv("ADDRESS", "")
	t.Setenv("ENABLE_CORS", "true")
	t.Setenv("FRONTEND_URL", "https://example.test")
	t.Setenv("ADMIN_PASSWORD", "test-admin-password")

	config, err := app.GetConfig()

	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	const want = "https://example.test"
	if config.Server.FrontendURL != want {
		t.Errorf("server FrontendURL = %q, want %q", config.Server.FrontendURL, want)
	}
	if config.Account.FrontendURL != want {
		t.Errorf("account FrontendURL = %q, want %q", config.Account.FrontendURL, want)
	}
	if config.Admin.FrontendURL != want {
		t.Errorf("admin FrontendURL = %q, want %q", config.Admin.FrontendURL, want)
	}
}
