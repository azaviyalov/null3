package journal_test

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/domain/journal"
	"github.com/azaviyalov/null3/backend/internal/testutil"
)

func TestServiceMoodRecordLifecycleAndOwnership(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newJournalTestEnvironment(t)
	owner := createJournalUser(t, environment, "owner")
	other := createJournalUser(t, environment, "other")

	record, err := environment.service.CreateMoodRecord(t.Context(), owner.ID, journal.MoodEditRecordRequest{
		Feeling: "calm",
		Emoji:   "🙂",
		Note:    "quiet morning",
	})
	if err != nil {
		t.Fatalf("CreateMoodRecord() error = %v", err)
	}
	if record.UserID != owner.ID {
		t.Errorf("CreateMoodRecord() user ID = %d, want %d", record.UserID, owner.ID)
	}
	if record.Feeling != "calm" || record.Emoji != "🙂" || record.Note != "quiet morning" {
		t.Error("CreateMoodRecord() did not preserve the requested fields")
	}

	if _, err := environment.service.GetMoodRecord(t.Context(), other.ID, record.ID); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("foreign GetMoodRecord() error = %v, want ErrItemNotFound", err)
	}
	if _, err := environment.service.UpdateMoodRecord(t.Context(), other.ID, record.ID, journal.MoodEditRecordRequest{Feeling: "changed"}); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("foreign UpdateMoodRecord() error = %v, want ErrItemNotFound", err)
	}
	if _, err := environment.service.DeleteMoodRecord(t.Context(), other.ID, record.ID); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("foreign DeleteMoodRecord() error = %v, want ErrItemNotFound", err)
	}

	updated, err := environment.service.UpdateMoodRecord(t.Context(), owner.ID, record.ID, journal.MoodEditRecordRequest{
		Feeling: "focused",
		Emoji:   "🎯",
		Note:    "deep work",
	})
	if err != nil {
		t.Fatalf("UpdateMoodRecord() error = %v", err)
	}
	if updated.Feeling != "focused" || updated.Emoji != "🎯" || updated.Note != "deep work" {
		t.Error("UpdateMoodRecord() did not preserve the requested fields")
	}

	deleted, err := environment.service.DeleteMoodRecord(t.Context(), owner.ID, record.ID)
	if err != nil {
		t.Fatalf("DeleteMoodRecord() error = %v", err)
	}
	if !deleted.DeletedAt.Valid {
		t.Fatal("DeleteMoodRecord() did not mark the record deleted")
	}
	if _, err := environment.service.DeleteMoodRecord(t.Context(), owner.ID, record.ID); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("repeated DeleteMoodRecord() error = %v, want ErrItemNotFound", err)
	}
	if _, err := environment.service.RestoreMoodRecord(t.Context(), other.ID, record.ID); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("foreign RestoreMoodRecord() error = %v, want ErrItemNotFound", err)
	}

	storedDeleted, err := environment.service.GetMoodRecord(t.Context(), owner.ID, record.ID)
	if err != nil {
		t.Fatalf("GetMoodRecord() deleted error = %v", err)
	}
	if !storedDeleted.DeletedAt.Valid {
		t.Fatal("GetMoodRecord() did not return the deleted state")
	}

	restored, err := environment.service.RestoreMoodRecord(t.Context(), owner.ID, record.ID)
	if err != nil {
		t.Fatalf("RestoreMoodRecord() error = %v", err)
	}
	if restored.DeletedAt.Valid {
		t.Fatal("RestoreMoodRecord() left the record deleted")
	}
	if _, err := environment.service.RestoreMoodRecord(t.Context(), owner.ID, record.ID); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("repeated RestoreMoodRecord() error = %v, want ErrItemNotFound", err)
	}
}

func TestServiceListMoodRecords(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newJournalTestEnvironment(t)
	owner := createJournalUser(t, environment, "owner")
	other := createJournalUser(t, environment, "other")
	baseTime := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	oldest := saveMoodRecord(t, environment, owner.ID, "oldest", baseTime)
	deleted := saveMoodRecord(t, environment, owner.ID, "deleted", baseTime.Add(time.Hour))
	newest := saveMoodRecord(t, environment, owner.ID, "newest", baseTime.Add(2*time.Hour))
	saveMoodRecord(t, environment, other.ID, "foreign", baseTime.Add(3*time.Hour))
	if _, err := environment.service.DeleteMoodRecord(t.Context(), owner.ID, deleted.ID); err != nil {
		t.Fatalf("delete mood record: %v", err)
	}

	firstPage, err := environment.service.ListMoodRecords(t.Context(), owner.ID, 1, 0, false)
	if err != nil {
		t.Fatalf("ListMoodRecords() error = %v", err)
	}
	if firstPage.TotalCount != 2 || !slices.Equal(moodRecordIDs(firstPage.Items), []uint{newest.ID}) {
		t.Fatalf("first active page IDs = %v total = %d, want [%d] and 2", moodRecordIDs(firstPage.Items), firstPage.TotalCount, newest.ID)
	}
	secondPage, err := environment.service.ListMoodRecords(t.Context(), owner.ID, 1, 1, false)
	if err != nil {
		t.Fatalf("ListMoodRecords() second page error = %v", err)
	}
	if secondPage.TotalCount != 2 || !slices.Equal(moodRecordIDs(secondPage.Items), []uint{oldest.ID}) {
		t.Fatalf("second active page IDs = %v total = %d, want [%d] and 2", moodRecordIDs(secondPage.Items), secondPage.TotalCount, oldest.ID)
	}
	deletedPage, err := environment.service.ListMoodRecords(t.Context(), owner.ID, 10, 0, true)
	if err != nil {
		t.Fatalf("ListMoodRecords() deleted error = %v", err)
	}
	if deletedPage.TotalCount != 1 || !slices.Equal(moodRecordIDs(deletedPage.Items), []uint{deleted.ID}) {
		t.Fatalf("deleted page IDs = %v total = %d, want [%d] and 1", moodRecordIDs(deletedPage.Items), deletedPage.TotalCount, deleted.ID)
	}
}
