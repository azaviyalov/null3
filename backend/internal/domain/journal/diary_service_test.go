package journal_test

import (
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/domain/journal"
	"github.com/azaviyalov/null3/backend/internal/testutil"
)

func TestServiceValidatesDiaryEntry(t *testing.T) {
	testutil.SkipIntegration(t)
	tests := []struct {
		name    string
		request journal.DiaryEditEntryRequest
	}{
		{name: "blank markdown", request: diaryRequest("  ", timePointer(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)))},
		{name: "missing occurred at", request: diaryRequest("entry", nil)},
		{name: "zero occurred at", request: diaryRequest("entry", timePointer(time.Time{}))},
		{name: "future occurred at", request: diaryRequest("entry", timePointer(time.Date(2100, time.January, 1, 0, 0, 0, 0, time.UTC)))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environment := newJournalTestEnvironment(t)
			owner := createJournalUser(t, environment, "owner")

			_, err := environment.service.CreateDiaryEntry(t.Context(), owner.ID, tt.request)

			if !errors.Is(err, core.ErrInvalidItem) {
				t.Fatalf("CreateDiaryEntry() error = %v, want ErrInvalidItem", err)
			}
			var count int64
			if err := environment.database.Model(&journal.DiaryEntry{}).Count(&count).Error; err != nil {
				t.Fatalf("count diary entries: %v", err)
			}
			if count != 0 {
				t.Fatal("invalid diary request created an entry")
			}
		})
	}
}

func TestServiceDiaryAssociationsAndOwnership(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newJournalTestEnvironment(t)
	owner := createJournalUser(t, environment, "owner")
	other := createJournalUser(t, environment, "other")
	createdAt := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	firstMood := saveMoodRecord(t, environment, owner.ID, "first", createdAt)
	secondMood := saveMoodRecord(t, environment, owner.ID, "second", createdAt.Add(time.Hour))
	deletedMood := saveMoodRecord(t, environment, owner.ID, "deleted", createdAt.Add(2*time.Hour))
	foreignMood := saveMoodRecord(t, environment, other.ID, "foreign", createdAt.Add(3*time.Hour))
	if _, err := environment.service.DeleteMoodRecord(t.Context(), owner.ID, deletedMood.ID); err != nil {
		t.Fatalf("delete mood record: %v", err)
	}

	localTime := time.Date(2026, time.January, 2, 12, 30, 0, 0, time.FixedZone("test", 3*60*60))
	markdown := fmt.Sprintf(
		"first [[mood:%d|First]] and /mood-records/%d; duplicate [[mood:%d]]; ignored `[[mood:%d]]`",
		firstMood.ID,
		secondMood.ID,
		firstMood.ID,
		foreignMood.ID,
	)
	entry, err := environment.service.CreateDiaryEntry(t.Context(), owner.ID, journal.DiaryEditEntryRequest{
		Title:      "  A day  ",
		Markdown:   "  " + markdown + "  ",
		OccurredAt: &localTime,
	})
	if err != nil {
		t.Fatalf("CreateDiaryEntry() error = %v", err)
	}
	if entry.UserID != owner.ID || entry.Title != "A day" || entry.Markdown != markdown || !entry.OccurredAt.Equal(localTime.UTC()) {
		t.Error("CreateDiaryEntry() did not normalize the request")
	}
	gotMoodIDs := moodRecordIDs(entry.MoodRecords)
	wantMoodIDs := []uint{secondMood.ID, firstMood.ID}
	if !slices.Equal(gotMoodIDs, wantMoodIDs) {
		t.Fatalf("created mood associations = %v, want %v", gotMoodIDs, wantMoodIDs)
	}

	invalidReferences := []struct {
		name     string
		markdown string
	}{
		{name: "foreign", markdown: fmt.Sprintf("[[mood:%d]]", foreignMood.ID)},
		{name: "deleted", markdown: fmt.Sprintf("[[mood:%d]]", deletedMood.ID)},
		{name: "missing", markdown: "[[mood:999999]]"},
	}
	for _, test := range invalidReferences {
		t.Run(test.name+" mood reference", func(t *testing.T) {
			_, err := environment.service.CreateDiaryEntry(t.Context(), owner.ID, diaryRequest(test.markdown, &localTime))
			if !errors.Is(err, core.ErrInvalidItem) {
				t.Fatalf("CreateDiaryEntry() error = %v, want ErrInvalidItem", err)
			}
		})
	}

	if _, err := environment.service.GetDiaryEntry(t.Context(), other.ID, entry.ID); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("foreign GetDiaryEntry() error = %v, want ErrItemNotFound", err)
	}
	if _, err := environment.service.UpdateDiaryEntry(t.Context(), other.ID, entry.ID, diaryRequest("changed", &localTime)); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("foreign UpdateDiaryEntry() error = %v, want ErrItemNotFound", err)
	}
	if _, err := environment.service.DeleteDiaryEntry(t.Context(), other.ID, entry.ID); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("foreign DeleteDiaryEntry() error = %v, want ErrItemNotFound", err)
	}

	updatedMarkdown := fmt.Sprintf("only [[mood:%d|Second]]", secondMood.ID)
	updated, err := environment.service.UpdateDiaryEntry(t.Context(), owner.ID, entry.ID, diaryRequest(updatedMarkdown, &localTime))
	if err != nil {
		t.Fatalf("UpdateDiaryEntry() error = %v", err)
	}
	if gotIDs := moodRecordIDs(updated.MoodRecords); !slices.Equal(gotIDs, []uint{secondMood.ID}) {
		t.Fatalf("updated mood associations = %v, want [%d]", gotIDs, secondMood.ID)
	}
	firstBacklinks, err := environment.service.GetMoodRecord(t.Context(), owner.ID, firstMood.ID)
	if err != nil {
		t.Fatalf("get first mood record: %v", err)
	}
	if len(firstBacklinks.DiaryEntries) != 0 {
		t.Fatal("association replacement left a stale backlink")
	}
	secondBacklinks, err := environment.service.GetMoodRecord(t.Context(), owner.ID, secondMood.ID)
	if err != nil {
		t.Fatalf("get second mood record: %v", err)
	}
	if len(secondBacklinks.DiaryEntries) != 1 || secondBacklinks.DiaryEntries[0].ID != entry.ID {
		t.Fatal("updated diary backlink is missing")
	}

	deletedEntry, err := environment.service.DeleteDiaryEntry(t.Context(), owner.ID, entry.ID)
	if err != nil {
		t.Fatalf("DeleteDiaryEntry() error = %v", err)
	}
	if !deletedEntry.DeletedAt.Valid {
		t.Fatal("DeleteDiaryEntry() did not mark the entry deleted")
	}
	if _, err := environment.service.DeleteDiaryEntry(t.Context(), owner.ID, entry.ID); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("repeated DeleteDiaryEntry() error = %v, want ErrItemNotFound", err)
	}
	if _, err := environment.service.RestoreDiaryEntry(t.Context(), other.ID, entry.ID); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("foreign RestoreDiaryEntry() error = %v, want ErrItemNotFound", err)
	}

	if _, err := environment.service.DeleteMoodRecord(t.Context(), owner.ID, secondMood.ID); err != nil {
		t.Fatalf("delete referenced mood record: %v", err)
	}
	if _, err := environment.service.RestoreDiaryEntry(t.Context(), owner.ID, entry.ID); !errors.Is(err, core.ErrInvalidItem) {
		t.Fatalf("RestoreDiaryEntry() with deleted mood error = %v, want ErrInvalidItem", err)
	}
	stillDeleted, err := environment.service.GetDiaryEntry(t.Context(), owner.ID, entry.ID)
	if err != nil {
		t.Fatalf("get rejected diary restore: %v", err)
	}
	if !stillDeleted.DeletedAt.Valid {
		t.Fatal("rejected diary restore changed its deleted state")
	}
	if _, err := environment.service.RestoreMoodRecord(t.Context(), owner.ID, secondMood.ID); err != nil {
		t.Fatalf("restore referenced mood record: %v", err)
	}
	restored, err := environment.service.RestoreDiaryEntry(t.Context(), owner.ID, entry.ID)
	if err != nil {
		t.Fatalf("RestoreDiaryEntry() error = %v", err)
	}
	if restored.DeletedAt.Valid || !slices.Equal(moodRecordIDs(restored.MoodRecords), []uint{secondMood.ID}) {
		t.Fatal("RestoreDiaryEntry() did not restore the entry and its mood association")
	}
}

func TestServiceListDiaryEntries(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newJournalTestEnvironment(t)
	owner := createJournalUser(t, environment, "owner")
	other := createJournalUser(t, environment, "other")
	oldestTime := time.Date(2026, time.January, 1, 8, 0, 0, 0, time.UTC)
	deletedTime := oldestTime.Add(time.Hour)
	newestTime := oldestTime.Add(2 * time.Hour)

	oldest, err := environment.service.CreateDiaryEntry(t.Context(), owner.ID, diaryRequest("oldest", &oldestTime))
	if err != nil {
		t.Fatalf("create oldest diary entry: %v", err)
	}
	deleted, err := environment.service.CreateDiaryEntry(t.Context(), owner.ID, diaryRequest("deleted", &deletedTime))
	if err != nil {
		t.Fatalf("create deleted diary entry: %v", err)
	}
	newest, err := environment.service.CreateDiaryEntry(t.Context(), owner.ID, diaryRequest("newest", &newestTime))
	if err != nil {
		t.Fatalf("create newest diary entry: %v", err)
	}
	if _, err := environment.service.CreateDiaryEntry(t.Context(), other.ID, diaryRequest("foreign", &newestTime)); err != nil {
		t.Fatalf("create foreign diary entry: %v", err)
	}
	if _, err := environment.service.DeleteDiaryEntry(t.Context(), owner.ID, deleted.ID); err != nil {
		t.Fatalf("delete diary entry: %v", err)
	}

	firstPage, err := environment.service.ListDiaryEntries(t.Context(), owner.ID, 1, 0, false)
	if err != nil {
		t.Fatalf("ListDiaryEntries() error = %v", err)
	}
	if firstPage.TotalCount != 2 || !slices.Equal(diaryEntryIDs(firstPage.Items), []uint{newest.ID}) {
		t.Fatalf("first active page IDs = %v total %d, want [%d] total 2", diaryEntryIDs(firstPage.Items), firstPage.TotalCount, newest.ID)
	}
	secondPage, err := environment.service.ListDiaryEntries(t.Context(), owner.ID, 1, 1, false)
	if err != nil {
		t.Fatalf("ListDiaryEntries() second page error = %v", err)
	}
	if secondPage.TotalCount != 2 || !slices.Equal(diaryEntryIDs(secondPage.Items), []uint{oldest.ID}) {
		t.Fatalf("second active page IDs = %v total %d, want [%d] total 2", diaryEntryIDs(secondPage.Items), secondPage.TotalCount, oldest.ID)
	}
	deletedPage, err := environment.service.ListDiaryEntries(t.Context(), owner.ID, 10, 0, true)
	if err != nil {
		t.Fatalf("ListDiaryEntries() deleted error = %v", err)
	}
	if deletedPage.TotalCount != 1 || !slices.Equal(diaryEntryIDs(deletedPage.Items), []uint{deleted.ID}) {
		t.Fatalf("deleted page IDs = %v total %d, want [%d] total 1", diaryEntryIDs(deletedPage.Items), deletedPage.TotalCount, deleted.ID)
	}
}

func TestServiceMoodRecordBacklinks(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newJournalTestEnvironment(t)
	owner := createJournalUser(t, environment, "owner")
	other := createJournalUser(t, environment, "other")
	mood := saveMoodRecord(t, environment, owner.ID, "linked", time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC))
	olderTime := time.Date(2026, time.January, 2, 8, 0, 0, 0, time.UTC)
	newerTime := olderTime.Add(time.Hour)
	older, err := environment.service.CreateDiaryEntry(t.Context(), owner.ID, diaryRequest(fmt.Sprintf("[[mood:%d]]", mood.ID), &olderTime))
	if err != nil {
		t.Fatalf("create older diary entry: %v", err)
	}
	newer, err := environment.service.CreateDiaryEntry(t.Context(), owner.ID, diaryRequest(fmt.Sprintf("/mood-records/%d", mood.ID), &newerTime))
	if err != nil {
		t.Fatalf("create newer diary entry: %v", err)
	}

	loaded, err := environment.service.GetMoodRecord(t.Context(), owner.ID, mood.ID)
	if err != nil {
		t.Fatalf("GetMoodRecord() error = %v", err)
	}
	gotBacklinkIDs := diaryEntryIDs(loaded.DiaryEntries)
	wantBacklinkIDs := []uint{newer.ID, older.ID}
	if !slices.Equal(gotBacklinkIDs, wantBacklinkIDs) {
		t.Fatalf("backlink IDs = %v, want %v", gotBacklinkIDs, wantBacklinkIDs)
	}
	if _, err := environment.service.GetMoodRecord(t.Context(), other.ID, mood.ID); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("foreign GetMoodRecord() error = %v, want ErrItemNotFound", err)
	}
	if _, err := environment.service.DeleteDiaryEntry(t.Context(), owner.ID, newer.ID); err != nil {
		t.Fatalf("delete newer diary entry: %v", err)
	}
	loaded, err = environment.service.GetMoodRecord(t.Context(), owner.ID, mood.ID)
	if err != nil {
		t.Fatalf("GetMoodRecord() after delete error = %v", err)
	}
	if gotIDs := diaryEntryIDs(loaded.DiaryEntries); !slices.Equal(gotIDs, []uint{older.ID}) {
		t.Fatalf("backlink IDs after delete = %v, want [%d]", gotIDs, older.ID)
	}
}

func diaryRequest(markdown string, occurredAt *time.Time) journal.DiaryEditEntryRequest {
	return journal.DiaryEditEntryRequest{Markdown: markdown, OccurredAt: occurredAt}
}

func timePointer(value time.Time) *time.Time {
	return &value
}

func diaryEntryIDs(entries []journal.DiaryEntry) []uint {
	ids := make([]uint, len(entries))
	for index, entry := range entries {
		ids[index] = entry.ID
	}
	return ids
}
