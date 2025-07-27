package mood

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/azaviyalov/null3/backend/internal/core"
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

func (r *Repository) GetEntry(filter *EntryFilter) (*Entry, error) {
	slog.Debug("getting entry", "filter", filter)

	var entry Entry

	if err := filter.Apply(r.DB).First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: entry not found", core.ErrItemNotFound)
		}
		return nil, fmt.Errorf("db error: %w", err)
	}
	return &entry, nil
}

func (r *Repository) ListEntries(filter *EntryFilter, limit, offset int) ([]Entry, error) {
	slog.Debug("listing entries", "filter", filter, "limit", limit, "offset", offset)

	var entries []Entry

	err := filter.Apply(r.DB).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *Repository) CountEntries(filter *EntryFilter) (int64, error) {
	slog.Debug("counting entries", "filter", filter)

	var count int64
	err := filter.Apply(r.DB.Model(&Entry{})).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("db error: %w", err)
	}
	return count, nil
}

func (r *Repository) SaveEntry(entry *Entry) (*Entry, error) {
	slog.Debug("saving entry", "entry", entry)

	if err := r.DB.Save(entry).Error; err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}
	return entry, nil
}

func (r *Repository) DeleteEntry(filter *EntryFilter) (*Entry, error) {
	slog.Debug("deleting entry", "filter", filter)

	var entry Entry

	// Check if the entry exists before deleting
	q := filter.Apply(r.DB).First(&entry)
	if err := q.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: entry not found", core.ErrItemNotFound)
		}
		return nil, fmt.Errorf("db error: %w", err)
	}

	if err := r.DB.Delete(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: entry not found", core.ErrItemNotFound)
		}
		return nil, fmt.Errorf("db error: %w", err)
	}
	return &entry, nil
}
