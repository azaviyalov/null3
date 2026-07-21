package session_test

import (
	"errors"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"github.com/golang-jwt/jwt/v5"
)

func TestServiceGeneratesScopedAccessTokens(t *testing.T) {
	config := session.Config{
		JWTSecret:     testJWTSecret,
		JWTExpiration: time.Hour,
	}
	service := session.NewService(nil, config)

	tests := []struct {
		name         string
		generate     func() (string, error)
		wantSubject  string
		wantScope    string
		wantLifetime time.Duration
	}{
		{
			name:         "user",
			generate:     func() (string, error) { return service.GenerateUserAccessToken(42) },
			wantSubject:  "42",
			wantScope:    "user",
			wantLifetime: config.JWTExpiration,
		},
		{
			name:         "admin",
			generate:     func() (string, error) { return service.GenerateAdminAccessToken(30 * time.Minute) },
			wantSubject:  "admin",
			wantScope:    "admin",
			wantLifetime: 30 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			tokenString, err := tt.generate()
			after := time.Now()
			if err != nil {
				t.Fatalf("generate token: %v", err)
			}

			token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (any, error) {
				return []byte(testJWTSecret), nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer("null3"), jwt.WithExpirationRequired())
			if err != nil {
				t.Fatalf("parse generated token: %v", err)
			}
			claims := token.Claims.(jwt.MapClaims)
			if claims["sub"] != tt.wantSubject || claims["scope"] != tt.wantScope || claims["iss"] != "null3" {
				t.Errorf("claims subject/scope/issuer = %v/%v/%v, want %s/%s/null3", claims["sub"], claims["scope"], claims["iss"], tt.wantSubject, tt.wantScope)
			}
			issuedAt, err := claims.GetIssuedAt()
			if err != nil || issuedAt == nil {
				t.Fatalf("issued-at claim error = %v", err)
			}
			expiresAt, err := claims.GetExpirationTime()
			if err != nil || expiresAt == nil {
				t.Fatalf("expiration claim error = %v", err)
			}
			if issuedAt.Time.Before(before.Add(-time.Second)) || issuedAt.Time.After(after.Add(time.Second)) {
				t.Errorf("issued-at = %v, want between %v and %v", issuedAt.Time, before, after)
			}
			if got := expiresAt.Time.Sub(issuedAt.Time); got != tt.wantLifetime {
				t.Errorf("token lifetime = %v, want %v", got, tt.wantLifetime)
			}
		})
	}
}

func TestServiceValidatesAccessTokenClaims(t *testing.T) {
	service := session.NewService(nil, session.Config{JWTSecret: testJWTSecret})
	now := time.Now()
	validUserClaims := jwt.MapClaims{
		"iss":   "null3",
		"sub":   "42",
		"scope": "user",
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
	}

	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{name: "malformed", token: "not-a-token", wantErr: session.ErrJWTInvalid},
		{name: "wrong signature", token: signClaims(t, jwt.SigningMethodHS256, "other-secret", validUserClaims), wantErr: session.ErrJWTInvalid},
		{name: "wrong algorithm", token: signClaims(t, jwt.SigningMethodHS384, testJWTSecret, validUserClaims), wantErr: session.ErrJWTInvalid},
		{name: "expired", token: signClaims(t, jwt.SigningMethodHS256, testJWTSecret, withClaim(validUserClaims, "exp", now.Add(-time.Minute).Unix())), wantErr: session.ErrJWTExpired},
		{name: "missing expiration", token: signClaims(t, jwt.SigningMethodHS256, testJWTSecret, withoutClaim(validUserClaims, "exp")), wantErr: session.ErrJWTInvalidClaims},
		{name: "wrong issuer", token: signClaims(t, jwt.SigningMethodHS256, testJWTSecret, withClaim(validUserClaims, "iss", "other")), wantErr: session.ErrJWTInvalidClaims},
		{name: "missing subject", token: signClaims(t, jwt.SigningMethodHS256, testJWTSecret, withoutClaim(validUserClaims, "sub")), wantErr: session.ErrJWTInvalidClaims},
		{name: "wrong scope", token: signClaims(t, jwt.SigningMethodHS256, testJWTSecret, withClaim(validUserClaims, "scope", "admin")), wantErr: session.ErrJWTInvalidClaims},
		{name: "non-numeric subject", token: signClaims(t, jwt.SigningMethodHS256, testJWTSecret, withClaim(validUserClaims, "sub", "user")), wantErr: session.ErrJWTInvalidClaims},
		{name: "zero subject", token: signClaims(t, jwt.SigningMethodHS256, testJWTSecret, withClaim(validUserClaims, "sub", "0")), wantErr: session.ErrJWTInvalidClaims},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := service.ParseUserAccessToken(tt.token)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ParseUserAccessToken() error = %v, want %v", err, tt.wantErr)
			}
			if userID != 0 {
				t.Errorf("ParseUserAccessToken() user ID = %d, want 0", userID)
			}
		})
	}

	validToken := signClaims(t, jwt.SigningMethodHS256, testJWTSecret, validUserClaims)
	userID, err := service.ParseUserAccessToken(validToken)
	if err != nil || userID != 42 {
		t.Fatalf("ParseUserAccessToken() = %d, %v; want 42, nil", userID, err)
	}
}

func TestServiceValidatesAdminAccessToken(t *testing.T) {
	service := session.NewService(nil, session.Config{JWTSecret: testJWTSecret})
	now := time.Now()
	validClaims := jwt.MapClaims{"iss": "null3", "sub": "admin", "scope": "admin", "exp": now.Add(time.Hour).Unix()}

	if err := service.ValidateAdminAccessToken(signClaims(t, jwt.SigningMethodHS256, testJWTSecret, validClaims)); err != nil {
		t.Fatalf("ValidateAdminAccessToken() error = %v", err)
	}

	tests := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{name: "wrong subject", claims: withClaim(validClaims, "sub", "42")},
		{name: "wrong scope", claims: withClaim(validClaims, "scope", "user")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateAdminAccessToken(signClaims(t, jwt.SigningMethodHS256, testJWTSecret, tt.claims))
			if !errors.Is(err, session.ErrJWTInvalidClaims) {
				t.Fatalf("ValidateAdminAccessToken() error = %v, want ErrJWTInvalidClaims", err)
			}
		})
	}
}

func TestServiceRefreshTokenLifecycle(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newSessionTestEnvironment(t)
	before := time.Now()

	first, err := environment.service.CreateRefreshToken(t.Context(), 41)
	if err != nil {
		t.Fatalf("CreateRefreshToken() error = %v", err)
	}
	second, err := environment.service.CreateRefreshToken(t.Context(), 42)
	if err != nil {
		t.Fatalf("CreateRefreshToken() second error = %v", err)
	}
	if first.Value == "" || second.Value == "" || first.Value == second.Value {
		t.Fatal("CreateRefreshToken() did not return distinct non-empty tokens")
	}
	if first.ExpiresAt.Before(before.Add(environment.config.RefreshTokenExpiration-time.Second)) || first.ExpiresAt.After(time.Now().Add(environment.config.RefreshTokenExpiration+time.Second)) {
		t.Errorf("refresh expiration = %v, want configured lifetime", first.ExpiresAt)
	}

	var stored session.RefreshToken
	if err := environment.database.First(&stored, first.ID).Error; err != nil {
		t.Fatalf("get stored refresh token: %v", err)
	}
	if stored.Value == first.Value || len(stored.Value) != 64 {
		t.Fatal("database contains a raw or malformed refresh-token value")
	}
	for _, tt := range []struct {
		token      string
		wantUserID uint
	}{
		{token: first.Value, wantUserID: 41},
		{token: second.Value, wantUserID: 42},
	} {
		found, err := environment.repository.GetRefreshToken(t.Context(), tt.token)
		if err != nil {
			t.Fatalf("GetRefreshToken() error = %v", err)
		}
		if found.UserID != tt.wantUserID {
			t.Errorf("GetRefreshToken() user ID = %d, want %d", found.UserID, tt.wantUserID)
		}
	}

	if err := environment.service.InvalidateRefreshToken(t.Context(), first.Value); err != nil {
		t.Fatalf("InvalidateRefreshToken() error = %v", err)
	}
	if _, err := environment.repository.GetRefreshToken(t.Context(), first.Value); err == nil {
		t.Fatal("invalidated refresh token is still available")
	}
	if err := environment.service.InvalidateRefreshToken(t.Context(), "unknown-token"); err != nil {
		t.Fatalf("InvalidateRefreshToken() unknown error = %v", err)
	}
}

func signClaims(t *testing.T, method jwt.SigningMethod, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(method, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}
	return token
}

func withClaim(claims jwt.MapClaims, name string, value any) jwt.MapClaims {
	copy := cloneClaims(claims)
	copy[name] = value
	return copy
}

func withoutClaim(claims jwt.MapClaims, name string) jwt.MapClaims {
	copy := cloneClaims(claims)
	delete(copy, name)
	return copy
}

func cloneClaims(claims jwt.MapClaims) jwt.MapClaims {
	copy := make(jwt.MapClaims, len(claims))
	for name, value := range claims {
		copy[name] = value
	}
	return copy
}
