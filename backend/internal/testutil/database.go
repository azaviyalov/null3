package testutil

import (
	"path/filepath"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/db"
	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/journal"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(t testing.TB, filename string) *gorm.DB {
	t.Helper()

	databaseURL := "file:" + filepath.Join(t.TempDir(), filename) + "?_fk=1"
	database, err := db.Connect(db.Config{DatabaseURL: databaseURL})
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	database.Logger = logger.Default.LogMode(logger.Silent)

	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("get SQL database: %v", err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("close SQL database: %v", err)
		}
	})

	if err := db.AutoMigrate(database,
		&journal.MoodRecord{},
		&journal.DiaryEntry{},
		&account.User{},
		&session.RefreshToken{},
		&account.PasswordResetToken{},
		&account.Invite{},
	); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	return database
}
