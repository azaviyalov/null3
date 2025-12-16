package mood

import (
	"context"
	"errors"
	"fmt"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
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
			// Not found is expected for some queries; log at Info so alerts aren't noisy
			logging.Info(ctx, "entry not found in GetEntry", "user_id", filter.UserID, "id", filter.ID)
			return nil, fmt.Errorf("%w: entry not found", core.ErrItemNotFound)
		}
		logging.Error(ctx, "db error in GetEntry", "error", err)
		return nil, fmt.Errorf("db error: %w", err)
	}
	logging.Info(ctx, "entry found in GetEntry", "user_id", entry.UserID, "id", entry.ID)
	return &entry, nil
}

func (r *Repository) ListEntries(ctx context.Context, filter *EntryFilter, limit, offset int) ([]Entry, error) {
	logging.Debug(ctx, "ListEntries called", "limit", limit, "offset", offset)

	var entries []Entry

	err := filter.Apply(r.DB.WithContext(ctx)).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries).Error
	if err != nil {
		logging.Error(ctx, "db error in ListEntries", "error", err, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("db error: %w", err)
	}

	logging.Info(ctx, "entries listed in ListEntries", "count", len(entries), "limit", limit, "offset", offset)
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
	logging.Info(ctx, "entries counted in CountEntries", "count", count)
	return count, nil
}

func (r *Repository) SaveEntry(ctx context.Context, entry *Entry) (*Entry, error) {
	if err := r.DB.WithContext(ctx).Save(entry).Error; err != nil {
		logging.Error(ctx, "db error in SaveEntry", "error", err, "user_id", entry.UserID)
		return nil, fmt.Errorf("db error: %w", err)
	}
	logging.Info(ctx, "entry saved in SaveEntry", "user_id", entry.UserID, "id", entry.ID)
	return entry, nil
}

func (r *Repository) DeleteEntry(ctx context.Context, filter *EntryFilter) (*Entry, error) {
	logging.Debug(ctx, "DeleteEntry called")

	var entry Entry

	// Check if the entry exists before deleting
	q := filter.Apply(r.DB.WithContext(ctx)).First(&entry)
	if err := q.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.Info(ctx, "entry not found in DeleteEntry", "user_id", filter.UserID, "id", filter.ID)
			return nil, fmt.Errorf("%w: entry not found", core.ErrItemNotFound)
		}
		logging.Error(ctx, "db error in DeleteEntry (find)", "error", err)
		return nil, fmt.Errorf("db error: %w", err)
	}

	if err := r.DB.WithContext(ctx).Delete(&entry).Error; err != nil {
		// log error for delete
		logging.Error(ctx, "db error in DeleteEntry (delete)", "error", err)
		return nil, fmt.Errorf("db error: %w", err)
	}
	logging.Info(ctx, "entry deleted in DeleteEntry", "user_id", entry.UserID, "id", entry.ID)
	return &entry, nil
}
