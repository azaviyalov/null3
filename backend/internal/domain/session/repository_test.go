package session_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"gorm.io/gorm"
)

func TestRepositoryRefreshTokenDeletion(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newSessionTestEnvironment(t)
	now := time.Now()
	tokens := []*session.RefreshToken{
		{UserID: 41, Value: "stored-one", CreatedAt: now, ExpiresAt: now.Add(-time.Hour)},
		{UserID: 41, Value: "stored-two", CreatedAt: now, ExpiresAt: now.Add(time.Hour)},
		{UserID: 42, Value: "stored-three", CreatedAt: now, ExpiresAt: now.Add(time.Hour)},
	}
	for _, token := range tokens {
		if _, err := environment.repository.SaveRefreshToken(t.Context(), token); err != nil {
			t.Fatalf("SaveRefreshToken() error = %v", err)
		}
	}

	if err := environment.repository.DeleteExpiredRefreshTokens(t.Context()); err != nil {
		t.Fatalf("DeleteExpiredRefreshTokens() error = %v", err)
	}
	assertRefreshTokenCount(t, environment, 2)

	if err := environment.repository.DeleteRefreshTokensByUser(t.Context(), 41); err != nil {
		t.Fatalf("DeleteRefreshTokensByUser() error = %v", err)
	}
	assertRefreshTokenCount(t, environment, 1)

	if err := environment.repository.DeleteRefreshToken(t.Context(), tokens[2]); err != nil {
		t.Fatalf("DeleteRefreshToken() error = %v", err)
	}
	assertRefreshTokenCount(t, environment, 0)
}

func TestRepositoryRefreshTokenErrors(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newSessionTestEnvironment(t)

	if _, err := environment.repository.GetRefreshToken(t.Context(), "unknown-token"); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("GetRefreshToken() error = %v, want core.ErrItemNotFound", err)
	}

	now := time.Now()
	duplicate := &session.RefreshToken{UserID: 41, Value: "duplicate", CreatedAt: now, ExpiresAt: now.Add(time.Hour)}
	if _, err := environment.repository.SaveRefreshToken(t.Context(), duplicate); err != nil {
		t.Fatalf("save first token: %v", err)
	}
	_, err := environment.repository.SaveRefreshToken(t.Context(), &session.RefreshToken{
		UserID: 42, Value: duplicate.Value, CreatedAt: now, ExpiresAt: now.Add(time.Hour),
	})
	if !errors.Is(err, gorm.ErrDuplicatedKey) {
		t.Fatalf("duplicate SaveRefreshToken() error = %v, want gorm.ErrDuplicatedKey", err)
	}

	canceledContext, cancel := context.WithCancel(t.Context())
	cancel()
	_, err = environment.repository.SaveRefreshToken(canceledContext, &session.RefreshToken{
		UserID: 43, Value: "canceled", CreatedAt: now, ExpiresAt: now.Add(time.Hour),
	})
	if !errors.Is(err, context.Canceled) || !strings.Contains(err.Error(), "save refresh token") {
		t.Fatalf("canceled SaveRefreshToken() error = %v, want contextual context.Canceled", err)
	}
}

func assertRefreshTokenCount(t *testing.T, environment *sessionTestEnvironment, want int64) {
	t.Helper()
	var count int64
	if err := environment.database.Model(&session.RefreshToken{}).Count(&count).Error; err != nil {
		t.Fatalf("count refresh tokens: %v", err)
	}
	if count != want {
		t.Fatalf("refresh token count = %d, want %d", count, want)
	}
}
