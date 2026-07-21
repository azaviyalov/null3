package account_test

import (
	"errors"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"gorm.io/gorm"
)

func TestRepositoryWithTx(t *testing.T) {
	testutil.SkipIntegration(t)
	tests := []struct {
		name             string
		login            string
		email            string
		transactionError error
		wantStored       bool
	}{
		{
			name:       "commit",
			login:      "committed",
			email:      "committed@example.test",
			wantStored: true,
		},
		{
			name:             "rollback",
			login:            "rolled-back",
			email:            "rolled-back@example.test",
			transactionError: errors.New("cancel transaction"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environment := newAccountTestEnvironment(t)

			err := environment.repository.WithTx(t.Context(), func(repository *account.Repository) error {
				_, err := repository.CreateUser(t.Context(), &account.User{
					Login:        tt.login,
					Email:        tt.email,
					PasswordHash: "test-hash",
				})
				if err != nil {
					return err
				}
				return tt.transactionError
			})

			if !errors.Is(err, tt.transactionError) {
				t.Fatalf("WithTx() error = %v, want %v", err, tt.transactionError)
			}
			_, err = environment.repository.GetUserByLogin(t.Context(), tt.login)
			if tt.wantStored && err != nil {
				t.Fatalf("get committed user: %v", err)
			}
			if !tt.wantStored && !errors.Is(err, core.ErrItemNotFound) {
				t.Fatalf("get rolled-back user error = %v, want ErrItemNotFound", err)
			}
		})
	}
}

func TestRepositoryTranslatesUniqueConstraints(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)
	createTestUser(t, environment, "existing", "existing@example.test")

	tests := []struct {
		name  string
		login string
		email string
	}{
		{name: "login", login: "existing", email: "different@example.test"},
		{name: "email", login: "different", email: "existing@example.test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := environment.repository.CreateUser(t.Context(), &account.User{
				Login:        tt.login,
				Email:        tt.email,
				PasswordHash: "test-hash",
			})

			if !errors.Is(err, gorm.ErrDuplicatedKey) {
				t.Fatalf("CreateUser() error = %v, want gorm.ErrDuplicatedKey", err)
			}
		})
	}
}
