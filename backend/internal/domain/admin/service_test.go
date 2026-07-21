package admin_test

import (
	"errors"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/admin"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/golang-jwt/jwt/v5"
)

const (
	testAdminPassword = "configured-admin-password"
	testJWTSecret     = "admin-test-signing-secret"
)

func TestServiceAuthenticate(t *testing.T) {
	tokenService := session.NewService(nil, session.Config{
		JWTSecret:     testJWTSecret,
		JWTExpiration: time.Hour,
	})
	service := admin.NewService(testAdminPassword, tokenService)

	t.Run("valid password", func(t *testing.T) {
		before := time.Now()
		token, err := service.Authenticate(testAdminPassword)
		after := time.Now()

		if err != nil {
			t.Fatalf("Authenticate() error = %v", err)
		}
		if token == "" {
			t.Fatal("Authenticate() returned an empty access token")
		}
		claims := parseAdminTokenClaims(t, token)
		if claims.Issuer != "null3" || claims.Subject != "admin" || claims.Scope != "admin" {
			t.Errorf("admin claims = issuer %q subject %q scope %q", claims.Issuer, claims.Subject, claims.Scope)
		}
		if claims.IssuedAt == nil || claims.ExpiresAt == nil {
			t.Fatal("admin token is missing issued-at or expiration")
		}
		if claims.IssuedAt.Time.Before(before.Add(-time.Second)) || claims.IssuedAt.Time.After(after.Add(time.Second)) {
			t.Errorf("issued-at = %v, want authentication time", claims.IssuedAt.Time)
		}
		wantExpiration := claims.IssuedAt.Time.Add(30 * time.Minute)
		if difference := claims.ExpiresAt.Time.Sub(wantExpiration); difference < -time.Second || difference > time.Second {
			t.Errorf("token lifetime difference = %v, want at most one second", difference)
		}
		if err := tokenService.ValidateAdminAccessToken(token); err != nil {
			t.Fatalf("ValidateAdminAccessToken() error = %v", err)
		}
	})

	t.Run("invalid passwords have one public error", func(t *testing.T) {
		tests := []struct {
			name     string
			password string
		}{
			{name: "empty", password: ""},
			{name: "incorrect", password: "incorrect-admin-password"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				token, err := service.Authenticate(tt.password)
				if !errors.Is(err, admin.ErrInvalidCredentials) {
					t.Fatalf("Authenticate() error = %v, want ErrInvalidCredentials", err)
				}
				if token != "" {
					t.Fatal("Authenticate() returned a token")
				}
			})
		}
	})
}

type adminTokenClaims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

func parseAdminTokenClaims(t *testing.T, tokenString string) *adminTokenClaims {
	t.Helper()

	claims := &adminTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			t.Fatalf("signing method = %s, want HS256", token.Method.Alg())
		}
		return []byte(testJWTSecret), nil
	})
	if err != nil {
		t.Fatalf("parse admin token: %v", err)
	}
	if !token.Valid {
		t.Fatal("admin token is invalid")
	}
	return claims
}
