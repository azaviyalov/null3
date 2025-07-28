package mood

import (
	"log/slog"

	"github.com/azaviyalov/null3/backend/internal/core"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ListEntries(userID uint, limit, offset int, deleted bool) (core.PaginatedResponse[Entry], error) {
	slog.Debug("ListEntries service called", "userID", userID, "limit", limit, "offset", offset, "deleted", deleted)
	filter := NewEntryFilter().WithUserID(userID)
	if deleted {
		filter = filter.WithDeletedMode(core.DeletedModeDeletedOnly)
	}
	slog.Debug("listing entries", "filter", filter)
	entries, err := s.repo.ListEntries(filter, limit, offset)
	if err != nil {
		slog.Error("failed to list entries", "error", err, "filter", filter)
		return core.PaginatedResponse[Entry]{}, err
	}
	totalCount, err := s.repo.CountEntries(filter)
	if err != nil {
		slog.Error("failed to count entries", "error", err, "filter", filter)
		return core.PaginatedResponse[Entry]{}, err
	}
	if entries == nil {
		entries = []Entry{}
	}
	slog.Info("successfully listed entries", "userID", userID, "count", len(entries), "limit", limit, "offset", offset, "deleted", deleted)
	return core.PaginatedResponse[Entry]{
		Items:      entries,
		TotalCount: totalCount,
	}, nil
}

func (s *Service) GetEntry(userID, id uint) (*Entry, error) {
	slog.Debug("GetEntry service called", "userID", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeAll)
	slog.Debug("getting entry", "filter", filter)
	entry, err := s.repo.GetEntry(filter)
	if err != nil {
		slog.Error("failed to get entry", "error", err, "filter", filter)
		return nil, err
	}
	slog.Info("successfully retrieved entry", "userID", userID, "id", id)
	return entry, nil
}

func (s *Service) CreateEntry(userID uint, req EditEntryRequest) (*Entry, error) {
	slog.Debug("CreateEntry service called", "userID", userID)
	entry, err := s.repo.SaveEntry(&Entry{
		UserID:  userID,
		Feeling: req.Feeling,
		Note:    req.Note,
	})
	if err != nil {
		slog.Error("failed to create entry", "error", err, "userID", userID)
		return nil, err
	}
	slog.Info("successfully created entry", "userID", userID, "id", entry.ID)
	return entry, nil
}

func (s *Service) UpdateEntry(userID, id uint, req EditEntryRequest) (*Entry, error) {
	slog.Debug("UpdateEntry service called", "userID", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id)
	slog.Debug("updating entry", "filter", filter)
	entry, err := s.repo.GetEntry(filter)
	if err != nil {
		slog.Error("failed to get entry for update", "error", err, "filter", filter)
		return nil, err
	}
	entry.Feeling = req.Feeling
	entry.Note = req.Note
	updatedEntry, err := s.repo.SaveEntry(entry)
	if err != nil {
		slog.Error("failed to update entry", "error", err, "userID", userID, "id", id)
		return nil, err
	}
	slog.Info("successfully updated entry", "userID", userID, "id", id)
	return updatedEntry, nil
}

func (s *Service) DeleteEntry(userID, id uint) (*Entry, error) {
	slog.Debug("DeleteEntry service called", "userID", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id)
	slog.Debug("deleting entry", "filter", filter)
	entry, err := s.repo.DeleteEntry(filter)
	if err != nil {
		slog.Error("failed to delete entry", "error", err, "filter", filter)
		return nil, err
	}
	slog.Info("successfully deleted entry", "userID", userID, "id", id)
	return entry, nil
}

func (s *Service) RestoreEntry(userID, id uint) (*Entry, error) {
	slog.Debug("RestoreEntry service called", "userID", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeDeletedOnly)
	slog.Debug("restoring entry", "filter", filter)
	entry, err := s.repo.GetEntry(filter)
	if err != nil {
		slog.Error("failed to get entry for restore", "error", err, "filter", filter)
		return nil, err
	}
	entry.DeletedAt.Valid = false
	restoredEntry, err := s.repo.SaveEntry(entry)
	if err != nil {
		slog.Error("failed to restore entry", "error", err, "userID", userID, "id", id)
		return nil, err
	}
	slog.Info("successfully restored entry", "userID", userID, "id", id)
	return restoredEntry, nil
}
