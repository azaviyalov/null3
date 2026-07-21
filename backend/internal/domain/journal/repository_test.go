package journal_test

import (
	"slices"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/journal"
	"github.com/azaviyalov/null3/backend/internal/testutil"
)

func TestRepositoryListMoodRecordsByIDsScopesOwnerAndDeletion(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newJournalTestEnvironment(t)
	owner := createJournalUser(t, environment, "owner")
	other := createJournalUser(t, environment, "other")
	createdAt := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	first := saveMoodRecord(t, environment, owner.ID, "first", createdAt)
	second := saveMoodRecord(t, environment, owner.ID, "second", createdAt.Add(time.Hour))
	deleted := saveMoodRecord(t, environment, owner.ID, "deleted", createdAt.Add(2*time.Hour))
	foreign := saveMoodRecord(t, environment, other.ID, "foreign", createdAt.Add(3*time.Hour))
	if _, err := environment.service.DeleteMoodRecord(t.Context(), owner.ID, deleted.ID); err != nil {
		t.Fatalf("delete mood record: %v", err)
	}

	records, err := environment.repository.ListMoodRecordsByIDs(t.Context(), owner.ID, []uint{
		foreign.ID,
		second.ID,
		deleted.ID,
		first.ID,
	})

	if err != nil {
		t.Fatalf("ListMoodRecordsByIDs() error = %v", err)
	}
	gotIDs := moodRecordIDs(records)
	wantIDs := []uint{first.ID, second.ID}
	if !slices.Equal(gotIDs, wantIDs) {
		t.Fatalf("ListMoodRecordsByIDs() IDs = %v, want %v", gotIDs, wantIDs)
	}

	empty, err := environment.repository.ListMoodRecordsByIDs(t.Context(), owner.ID, nil)
	if err != nil {
		t.Fatalf("ListMoodRecordsByIDs(nil) error = %v", err)
	}
	if empty == nil || len(empty) != 0 {
		t.Fatalf("ListMoodRecordsByIDs(nil) = %v, want non-nil empty slice", empty)
	}
}

func moodRecordIDs(records []journal.MoodRecord) []uint {
	ids := make([]uint, len(records))
	for index, record := range records {
		ids[index] = record.ID
	}
	return ids
}
