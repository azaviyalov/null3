package journal_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/journal"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"gorm.io/gorm"
)

type journalTestEnvironment struct {
	database   *gorm.DB
	repository *journal.Repository
	service    *journal.Service
}

func newJournalTestEnvironment(t *testing.T) *journalTestEnvironment {
	t.Helper()

	database := testutil.NewDatabase(t, "journal.sqlite")

	repository := journal.NewRepository(database)
	return &journalTestEnvironment{
		database:   database,
		repository: repository,
		service:    journal.NewService(repository),
	}
}

func createJournalUser(t *testing.T, environment *journalTestEnvironment, name string) *account.User {
	t.Helper()

	user := &account.User{
		Login:        name,
		Email:        fmt.Sprintf("%s@example.test", name),
		PasswordHash: "test-password-hash",
	}
	if err := environment.database.Create(user).Error; err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return user
}

func saveMoodRecord(t *testing.T, environment *journalTestEnvironment, userID uint, feeling string, createdAt time.Time) *journal.MoodRecord {
	t.Helper()

	record, err := environment.repository.SaveMoodRecord(t.Context(), &journal.MoodRecord{
		UserID:    userID,
		Feeling:   feeling,
		CreatedAt: createdAt,
	})
	if err != nil {
		t.Fatalf("save mood record: %v", err)
	}
	return record
}
