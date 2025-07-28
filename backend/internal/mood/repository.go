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
	slog.Debug("GetEntry called", "filter", filter)

	var entry Entry

	if err := filter.Apply(r.DB).First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("entry not found in GetEntry", "filter", filter)
			return nil, fmt.Errorf("%w: entry not found", core.ErrItemNotFound)
		}
		slog.Error("db error in GetEntry", "error", err, "filter", filter)
		return nil, fmt.Errorf("db error: %w", err)
	}
	slog.Info("entry found in GetEntry", "userID", entry.UserID, "id", entry.ID)
	return &entry, nil
}

func (r *Repository) ListEntries(filter *EntryFilter, limit, offset int) ([]Entry, error) {
	slog.Debug("ListEntries called", "filter", filter, "limit", limit, "offset", offset)

	var entries []Entry

	err := filter.Apply(r.DB).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries).Error
	if err != nil {
		slog.Error("db error in ListEntries", "error", err, "filter", filter, "limit", limit, "offset", offset)
		return nil, err
	}

	slog.Info("entries listed in ListEntries", "filter", filter, "count", len(entries), "limit", limit, "offset", offset)
	return entries, nil
}

func (r *Repository) CountEntries(filter *EntryFilter) (int64, error) {
	slog.Debug("CountEntries called", "filter", filter)

	var count int64
	err := filter.Apply(r.DB.Model(&Entry{})).Count(&count).Error
	if err != nil {
		slog.Error("db error in CountEntries", "error", err, "filter", filter)
		return 0, fmt.Errorf("db error: %w", err)
	}
	slog.Info("entries counted in CountEntries", "filter", filter, "count", count)
	return count, nil
}

func (r *Repository) SaveEntry(entry *Entry) (*Entry, error) {
	if err := r.DB.Save(entry).Error; err != nil {
		slog.Error("db error in SaveEntry", "error", err, "userID", entry.UserID)
		return nil, fmt.Errorf("db error: %w", err)
	}
	slog.Info("entry saved in SaveEntry", "userID", entry.UserID, "id", entry.ID)
	return entry, nil
}

func (r *Repository) DeleteEntry(filter *EntryFilter) (*Entry, error) {
	slog.Debug("DeleteEntry called", "filter", filter)

	var entry Entry

	// Check if the entry exists before deleting
	q := filter.Apply(r.DB).First(&entry)
	if err := q.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("entry not found in DeleteEntry", "filter", filter)
			return nil, fmt.Errorf("%w: entry not found", core.ErrItemNotFound)
		}
		slog.Error("db error in DeleteEntry (find)", "error", err, "filter", filter)
		return nil, fmt.Errorf("db error: %w", err)
	}

	if err := r.DB.Delete(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("entry not found in DeleteEntry (delete)", "filter", filter)
		}
		// log error for delete
		slog.Error("db error in DeleteEntry (delete)", "error", err, "filter", filter)
		return nil, fmt.Errorf("db error: %w", err)
	}
	slog.Info("entry deleted in DeleteEntry", "userID", entry.UserID, "id", entry.ID)
	return &entry, nil
}
