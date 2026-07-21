package session_test

import (
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"gorm.io/gorm"
)

const testJWTSecret = "session-test-signing-secret"

type sessionTestEnvironment struct {
	database   *gorm.DB
	repository *session.Repository
	service    *session.Service
	config     session.Config
}

func newSessionTestEnvironment(t *testing.T) *sessionTestEnvironment {
	t.Helper()

	database := testutil.NewDatabase(t, "session.sqlite")

	config := session.Config{
		JWTSecret:              testJWTSecret,
		JWTExpiration:          time.Hour,
		RefreshTokenExpiration: 7 * 24 * time.Hour,
	}
	repository := session.NewRepository(database)

	return &sessionTestEnvironment{
		database:   database,
		repository: repository,
		service:    session.NewService(repository, config),
		config:     config,
	}
}
