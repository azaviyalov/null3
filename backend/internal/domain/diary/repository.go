package diary

import (
	"context"
	"errors"
	"fmt"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	mooddomain "github.com/azaviyalov/null3/backend/internal/domain/mood"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) GetEntry(ctx context.Context, filter *EntryFilter) (*Entry, error) {
	logging.Debug(ctx, "GetEntry called")

	var entry Entry
	if err := filter.Apply(r.DB.WithContext(ctx)).First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.Info(ctx, "diary entry not found in GetEntry", "user_id", filter.UserID, "id", filter.ID)
			return nil, fmt.Errorf("%w: diary entry not found", core.ErrItemNotFound)
		}

		logging.Error(ctx, "db error in GetEntry", "error", err)
		return nil, fmt.Errorf("db error: %w", err)
	}

	logging.Info(ctx, "diary entry found in GetEntry", "user_id", entry.UserID, "id", entry.ID)
	return &entry, nil
}

func (r *Repository) ListEntries(ctx context.Context, filter *EntryFilter, limit, offset int) ([]Entry, error) {
	logging.Debug(ctx, "ListEntries called", "limit", limit, "offset", offset)

	var entries []Entry
	err := filter.Apply(r.DB.WithContext(ctx)).
		Order("occurred_at DESC").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries).Error
	if err != nil {
		logging.Error(ctx, "db error in ListEntries", "error", err, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("db error: %w", err)
	}

	logging.Info(ctx, "diary entries listed in ListEntries", "count", len(entries), "limit", limit, "offset", offset)
	return entries, nil
}

func (r *Repository) CountEntries(ctx context.Context, filter *EntryFilter) (int64, error) {
	logging.Debug(ctx, "CountEntries called")

	var count int64
	err := filter.Apply(r.DB.WithContext(ctx).Model(&Entry{})).Count(&count).Error
	if err != nil {
		logging.Error(ctx, "db error in CountEntries", "error", err)
		return 0, fmt.Errorf("db error: %w", err)
	}

	logging.Info(ctx, "diary entries counted in CountEntries", "count", count)
	return count, nil
}

func (r *Repository) SaveEntry(ctx context.Context, entry *Entry, referencedMoodEntryIDs []uint) (*Entry, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.validateMoodEntryIDs(tx, entry.UserID, referencedMoodEntryIDs); err != nil {
			return err
		}

		if err := tx.Save(entry).Error; err != nil {
			return fmt.Errorf("save diary entry: %w", err)
		}

		if err := tx.Where("diary_entry_id = ?", entry.ID).Delete(&MoodEntryLink{}).Error; err != nil {
			return fmt.Errorf("clear diary entry mood entry links: %w", err)
		}

		if len(referencedMoodEntryIDs) == 0 {
			return nil
		}

		references := make([]MoodEntryLink, 0, len(referencedMoodEntryIDs))
		for _, moodEntryID := range referencedMoodEntryIDs {
			references = append(references, MoodEntryLink{
				DiaryEntryID: entry.ID,
				MoodEntryID:  moodEntryID,
			})
		}

		if err := tx.Create(&references).Error; err != nil {
			return fmt.Errorf("create diary entry mood entry links: %w", err)
		}

		return nil
	})
	if err != nil {
		logging.Error(ctx, "db error in SaveEntry", "error", err, "user_id", entry.UserID, "id", entry.ID)
		return nil, err
	}

	logging.Info(ctx, "diary entry saved in SaveEntry", "user_id", entry.UserID, "id", entry.ID)
	return entry, nil
}

func (r *Repository) DeleteEntry(ctx context.Context, filter *EntryFilter) (*Entry, error) {
	logging.Debug(ctx, "DeleteEntry called")

	var entry Entry
	q := filter.Apply(r.DB.WithContext(ctx)).First(&entry)
	if err := q.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.Info(ctx, "diary entry not found in DeleteEntry", "user_id", filter.UserID, "id", filter.ID)
			return nil, fmt.Errorf("%w: diary entry not found", core.ErrItemNotFound)
		}

		logging.Error(ctx, "db error in DeleteEntry (find)", "error", err)
		return nil, fmt.Errorf("db error: %w", err)
	}

	if err := r.DB.WithContext(ctx).Delete(&entry).Error; err != nil {
		logging.Error(ctx, "db error in DeleteEntry (delete)", "error", err)
		return nil, fmt.Errorf("db error: %w", err)
	}

	logging.Info(ctx, "diary entry deleted in DeleteEntry", "user_id", entry.UserID, "id", entry.ID)
	return &entry, nil
}

func (r *Repository) ListReferencedMoodEntries(ctx context.Context, userID, diaryEntryID uint) ([]ReferencedMoodEntry, error) {
	logging.Debug(ctx, "ListReferencedMoodEntries called", "user_id", userID, "diary_entry_id", diaryEntryID)

	var moodEntries []mooddomain.Entry
	moodEntryLinkTable := (&MoodEntryLink{}).TableName()
	err := r.DB.WithContext(ctx).
		Unscoped().
		Model(&mooddomain.Entry{}).
		Joins(fmt.Sprintf("JOIN %s ON %s.mood_entry_id = mood_entries.id", moodEntryLinkTable, moodEntryLinkTable)).
		Where(fmt.Sprintf("%s.diary_entry_id = ?", moodEntryLinkTable), diaryEntryID).
		Where("mood_entries.user_id = ?", userID).
		Order("mood_entries.created_at DESC").
		Find(&moodEntries).Error
	if err != nil {
		logging.Error(ctx, "db error in ListReferencedMoodEntries", "error", err, "user_id", userID, "diary_entry_id", diaryEntryID)
		return nil, fmt.Errorf("db error: %w", err)
	}

	references := make([]ReferencedMoodEntry, 0, len(moodEntries))
	for _, moodEntry := range moodEntries {
		references = append(references, NewReferencedMoodEntry(moodEntry))
	}

	logging.Info(ctx, "referenced mood entries listed", "count", len(references), "user_id", userID, "diary_entry_id", diaryEntryID)
	return references, nil
}

func (r *Repository) ListDiaryEntryLinks(ctx context.Context, userID, moodEntryID uint) ([]mooddomain.DiaryEntryLink, error) {
	logging.Debug(ctx, "ListDiaryEntryLinks called", "user_id", userID, "mood_entry_id", moodEntryID)

	var entries []Entry
	moodEntryLinkTable := (&MoodEntryLink{}).TableName()
	err := r.DB.WithContext(ctx).
		Model(&Entry{}).
		Joins(fmt.Sprintf("JOIN %s ON %s.diary_entry_id = diary_entries.id", moodEntryLinkTable, moodEntryLinkTable)).
		Where("diary_entries.user_id = ?", userID).
		Where(fmt.Sprintf("%s.mood_entry_id = ?", moodEntryLinkTable), moodEntryID).
		Order("diary_entries.occurred_at DESC").
		Order("diary_entries.created_at DESC").
		Find(&entries).Error
	if err != nil {
		logging.Error(ctx, "db error in ListDiaryEntryLinks", "error", err, "user_id", userID, "mood_entry_id", moodEntryID)
		return nil, fmt.Errorf("db error: %w", err)
	}

	diaryEntryLinks := make([]mooddomain.DiaryEntryLink, 0, len(entries))
	for _, entry := range entries {
		diaryEntryLinks = append(diaryEntryLinks, mooddomain.DiaryEntryLink{
			ID:         entry.ID,
			Title:      entry.Title,
			Preview:    MarkdownPreview(entry.Markdown),
			OccurredAt: entry.OccurredAt,
			CreatedAt:  entry.CreatedAt,
			UpdatedAt:  entry.UpdatedAt,
		})
	}

	logging.Info(ctx, "diary entry links listed", "count", len(diaryEntryLinks), "user_id", userID, "mood_entry_id", moodEntryID)
	return diaryEntryLinks, nil
}

func (r *Repository) validateMoodEntryIDs(tx *gorm.DB, userID uint, moodEntryIDs []uint) error {
	if len(moodEntryIDs) == 0 {
		return nil
	}

	var count int64
	err := tx.Unscoped().
		Model(&mooddomain.Entry{}).
		Where("user_id = ?", userID).
		Where("id IN ?", moodEntryIDs).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("validate referenced mood entries: %w", err)
	}

	if count != int64(len(moodEntryIDs)) {
		return fmt.Errorf("%w: one or more referenced mood entries do not exist", core.ErrInvalidItem)
	}

	return nil
}
