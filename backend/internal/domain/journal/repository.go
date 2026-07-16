package journal

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/azaviyalov/null3/backend/internal/core"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetMoodRecord(ctx context.Context, filter *MoodRecordFilter) (*MoodRecord, error) {
	var entry MoodRecord
	query := filter.Apply(r.db.WithContext(ctx)).
		Preload("DiaryEntries", func(db *gorm.DB) *gorm.DB {
			return db.Order("occurred_at DESC").Order("created_at DESC")
		})
	if err := query.First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: mood record not found", core.ErrItemNotFound)
		}
		return nil, fmt.Errorf("get mood record: %w", err)
	}
	return &entry, nil
}

func (r *Repository) ListMoodRecords(ctx context.Context, filter *MoodRecordFilter, limit, offset int) ([]MoodRecord, error) {
	var entries []MoodRecord
	err := filter.Apply(r.db.WithContext(ctx)).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("list mood records: %w", err)
	}
	return entries, nil
}

func (r *Repository) CountMoodRecords(ctx context.Context, filter *MoodRecordFilter) (int64, error) {
	var count int64
	err := filter.Apply(r.db.WithContext(ctx).Model(&MoodRecord{})).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count mood records: %w", err)
	}
	return count, nil
}

func (r *Repository) SaveMoodRecord(ctx context.Context, entry *MoodRecord) (*MoodRecord, error) {
	if err := r.db.WithContext(ctx).Save(entry).Error; err != nil {
		return nil, fmt.Errorf("save mood record: %w", err)
	}
	return entry, nil
}

func (r *Repository) DeleteMoodRecord(ctx context.Context, filter *MoodRecordFilter) (*MoodRecord, error) {
	var entry MoodRecord
	q := filter.Apply(r.db.WithContext(ctx)).First(&entry)
	if err := q.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: mood record not found", core.ErrItemNotFound)
		}
		return nil, fmt.Errorf("find mood record to delete: %w", err)
	}
	if err := r.db.WithContext(ctx).Delete(&entry).Error; err != nil {
		return nil, fmt.Errorf("delete mood record: %w", err)
	}
	return &entry, nil
}

func (r *Repository) GetDiaryEntry(ctx context.Context, filter *DiaryEntryFilter) (*DiaryEntry, error) {
	var entry DiaryEntry
	query := filter.Apply(r.db.WithContext(ctx)).
		Preload("MoodRecords", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		})
	if err := query.First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: diary entry not found", core.ErrItemNotFound)
		}
		return nil, fmt.Errorf("get diary entry: %w", err)
	}
	return &entry, nil
}

func (r *Repository) ListDiaryEntries(ctx context.Context, filter *DiaryEntryFilter, limit, offset int) ([]DiaryEntry, error) {
	var entries []DiaryEntry
	err := filter.Apply(r.db.WithContext(ctx)).
		Order("occurred_at DESC").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("list diary entries: %w", err)
	}
	return entries, nil
}

func (r *Repository) CountDiaryEntries(ctx context.Context, filter *DiaryEntryFilter) (int64, error) {
	var count int64
	err := filter.Apply(r.db.WithContext(ctx).Model(&DiaryEntry{})).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count diary entries: %w", err)
	}
	return count, nil
}

func (r *Repository) SaveDiaryEntry(ctx context.Context, entry *DiaryEntry) (*DiaryEntry, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("MoodRecords").Save(entry).Error; err != nil {
			return fmt.Errorf("save diary entry: %w", err)
		}
		if err := tx.Model(entry).Association("MoodRecords").Replace(entry.MoodRecords); err != nil {
			return fmt.Errorf("replace diary mood links: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	filter := NewDiaryEntryFilter().
		WithUserID(entry.UserID).
		WithID(entry.ID).
		WithDeletedMode(core.DeletedModeAll)

	updatedEntry, err := r.GetDiaryEntry(ctx, filter)
	if err != nil {
		return nil, err
	}

	return updatedEntry, nil
}

func (r *Repository) ListMoodRecordsByIDs(ctx context.Context, userID uint, ids []uint) ([]MoodRecord, error) {
	if len(ids) == 0 {
		return []MoodRecord{}, nil
	}

	var entries []MoodRecord
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("id IN ?", ids).
		Find(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("list referenced mood records: %w", err)
	}

	slices.SortFunc(entries, func(a, b MoodRecord) int {
		switch {
		case a.ID < b.ID:
			return -1
		case a.ID > b.ID:
			return 1
		default:
			return 0
		}
	})

	return entries, nil
}

func (r *Repository) DeleteDiaryEntry(ctx context.Context, filter *DiaryEntryFilter) (*DiaryEntry, error) {
	var entry DiaryEntry
	q := filter.Apply(r.db.WithContext(ctx)).First(&entry)
	if err := q.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: diary entry not found", core.ErrItemNotFound)
		}
		return nil, fmt.Errorf("find diary entry to delete: %w", err)
	}
	if err := r.db.WithContext(ctx).Delete(&entry).Error; err != nil {
		return nil, fmt.Errorf("delete diary entry: %w", err)
	}
	return &entry, nil
}
