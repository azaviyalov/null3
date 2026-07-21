package session_test

import (
	"strings"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/session"
)

func TestGetConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		setSessionEnvironment(t)
		t.Setenv("JWT_SECRET", "  configured secret  ")

		config, err := session.GetConfig()
		if err != nil {
			t.Fatalf("GetConfig() error = %v", err)
		}
		if config.JWTSecret != "  configured secret  " {
			t.Error("GetConfig() changed JWT_SECRET")
		}
		if config.JWTExpiration != 24*time.Hour {
			t.Errorf("JWTExpiration = %v, want %v", config.JWTExpiration, 24*time.Hour)
		}
		if config.RefreshTokenExpiration != 7*24*time.Hour {
			t.Errorf("RefreshTokenExpiration = %v, want %v", config.RefreshTokenExpiration, 7*24*time.Hour)
		}
		if config.SecureCookies {
			t.Error("SecureCookies = true, want false")
		}
	})

	t.Run("overrides", func(t *testing.T) {
		setSessionEnvironment(t)
		t.Setenv("JWT_SECRET", testJWTSecret)
		t.Setenv("JWT_EXPIRATION", "90m")
		t.Setenv("REFRESH_TOKEN_EXPIRATION", "48h")
		t.Setenv("SECURE_COOKIES", "true")

		config, err := session.GetConfig()
		if err != nil {
			t.Fatalf("GetConfig() error = %v", err)
		}
		if config.JWTExpiration != 90*time.Minute {
			t.Errorf("JWTExpiration = %v, want %v", config.JWTExpiration, 90*time.Minute)
		}
		if config.RefreshTokenExpiration != 48*time.Hour {
			t.Errorf("RefreshTokenExpiration = %v, want %v", config.RefreshTokenExpiration, 48*time.Hour)
		}
		if !config.SecureCookies {
			t.Error("SecureCookies = false, want true")
		}
	})

	tests := []struct {
		name      string
		variable  string
		value     string
		wantError string
	}{
		{name: "missing secret", variable: "JWT_SECRET", wantError: "JWT_SECRET must be set"},
		{name: "blank secret", variable: "JWT_SECRET", value: "  ", wantError: "JWT_SECRET must be set"},
		{name: "invalid JWT expiration", variable: "JWT_EXPIRATION", value: "later", wantError: "parse JWT_EXPIRATION"},
		{name: "non-positive JWT expiration", variable: "JWT_EXPIRATION", value: "0s", wantError: "JWT_EXPIRATION must be a positive duration"},
		{name: "invalid refresh expiration", variable: "REFRESH_TOKEN_EXPIRATION", value: "later", wantError: "parse REFRESH_TOKEN_EXPIRATION"},
		{name: "non-positive refresh expiration", variable: "REFRESH_TOKEN_EXPIRATION", value: "-1s", wantError: "REFRESH_TOKEN_EXPIRATION must be a positive duration"},
		{name: "invalid secure cookies", variable: "SECURE_COOKIES", value: "sometimes", wantError: "parse SECURE_COOKIES"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setSessionEnvironment(t)
			t.Setenv("JWT_SECRET", testJWTSecret)
			t.Setenv(tt.variable, tt.value)

			_, err := session.GetConfig()
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("GetConfig() error = %v, want text %q", err, tt.wantError)
			}
		})
	}
}

func setSessionEnvironment(t *testing.T) {
	t.Helper()
	for _, name := range []string{"JWT_SECRET", "JWT_EXPIRATION", "REFRESH_TOKEN_EXPIRATION", "SECURE_COOKIES"} {
		t.Setenv(name, "")
	}
}
