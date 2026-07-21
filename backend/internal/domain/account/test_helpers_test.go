package account_test

import (
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const testPassword = "correct-password"

type accountTestEnvironment struct {
	database       *gorm.DB
	repository     *account.Repository
	service        *account.Service
	sessionService *session.Service
	sessionConfig  session.Config
	accountConfig  account.Config
}

func newAccountTestEnvironment(t *testing.T) *accountTestEnvironment {
	t.Helper()

	database := testutil.NewDatabase(t, "account.sqlite")

	sessionConfig := session.Config{
		JWTSecret:              "account-test-signing-secret",
		JWTExpiration:          time.Hour,
		RefreshTokenExpiration: 7 * 24 * time.Hour,
	}
	sessionService := session.NewService(session.NewRepository(database), sessionConfig)
	repository := account.NewRepository(database)
	accountConfig := account.Config{
		PasswordResetTokenExpiration: time.Hour,
		FrontendURL:                  "https://journal.example",
	}

	return &accountTestEnvironment{
		database:       database,
		repository:     repository,
		sessionService: sessionService,
		sessionConfig:  sessionConfig,
		accountConfig:  accountConfig,
		service:        account.NewService(repository, sessionService, accountConfig),
	}
}

func createTestUser(t *testing.T, environment *accountTestEnvironment, login, email string) *account.User {
	t.Helper()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash test password: %v", err)
	}
	user, err := environment.repository.CreateUser(t.Context(), &account.User{
		Login:        login,
		Email:        email,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return user
}

func assertUserResponse(t *testing.T, got *account.UserResponse, want *account.User) {
	t.Helper()
	if got == nil {
		t.Fatal("user response is nil")
	}
	if got.ID != want.ID || got.Login != want.Login || got.Email != want.Email {
		t.Errorf(
			"user response = ID %d login %q email %q; want ID %d login %q email %q",
			got.ID, got.Login, got.Email, want.ID, want.Login, want.Email,
		)
	}
}
