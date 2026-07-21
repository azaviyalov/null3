package db_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/db"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"gorm.io/gorm"
)

func TestGetConfig(t *testing.T) {
	t.Run("default database URL", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "")

		got := db.GetConfig()

		const want = "file:null3.db?_fk=1"
		if got.DatabaseURL != want {
			t.Fatalf("DatabaseURL = %q, want %q", got.DatabaseURL, want)
		}
	})

	t.Run("environment override", func(t *testing.T) {
		const want = "file:custom.db?_fk=1"
		t.Setenv("DATABASE_URL", want)

		got := db.GetConfig()

		if got.DatabaseURL != want {
			t.Fatalf("DatabaseURL = %q, want %q", got.DatabaseURL, want)
		}
	})
}

func TestConnectWrapsOpenError(t *testing.T) {
	testutil.SkipIntegration(t)
	databaseURL := "file:" + filepath.Join(t.TempDir(), "missing", "database.sqlite") + "?mode=ro"

	_, err := db.Connect(db.Config{DatabaseURL: databaseURL})

	if err == nil {
		t.Fatal("Connect() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "connect to database") {
		t.Fatalf("Connect() error = %q, want connection context", err)
	}
	if errors.Unwrap(err) == nil {
		t.Fatal("Connect() error does not retain its cause")
	}
}

func TestAutoMigrate(t *testing.T) {
	testutil.SkipIntegration(t)
	t.Run("creates model table", func(t *testing.T) {
		database := openTestDatabase(t)

		if err := db.AutoMigrate(database, &migrationRecord{}); err != nil {
			t.Fatalf("AutoMigrate() error = %v", err)
		}
		if !database.Migrator().HasTable(&migrationRecord{}) {
			t.Fatal("AutoMigrate() did not create the model table")
		}
	})

	t.Run("wraps migration error", func(t *testing.T) {
		database := openTestDatabase(t)
		sqlDB, err := database.DB()
		if err != nil {
			t.Fatalf("get SQL database: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Fatalf("close SQL database: %v", err)
		}

		err = db.AutoMigrate(database, &migrationRecord{})

		if err == nil {
			t.Fatal("AutoMigrate() error = nil, want an error")
		}
		if !strings.Contains(err.Error(), "migrate database") {
			t.Fatalf("AutoMigrate() error = %q, want migration context", err)
		}
		if errors.Unwrap(err) == nil {
			t.Fatal("AutoMigrate() error does not retain its cause")
		}
	})
}

type migrationRecord struct {
	ID   uint
	Name string
}

func openTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()

	databaseURL := filepath.Join(t.TempDir(), "database.sqlite")
	database, err := db.Connect(db.Config{DatabaseURL: databaseURL})
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("get SQL database: %v", err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("close SQL database: %v", err)
		}
	})
	return database
}
