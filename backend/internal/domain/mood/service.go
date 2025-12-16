package mood

import (
	"context"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ListEntries(ctx context.Context, userID uint, limit, offset int, deleted bool) (core.PaginatedResponse[Entry], error) {
	logging.Debug(ctx, "ListEntries service called", "user_id", userID, "limit", limit, "offset", offset, "deleted", deleted)
	filter := NewEntryFilter().WithUserID(userID)
	if deleted {
		filter = filter.WithDeletedMode(core.DeletedModeDeletedOnly)
	}
	logging.Debug(ctx, "listing entries", "filter", filter)
	entries, err := s.repo.ListEntries(ctx, filter, limit, offset)
	if err != nil {
		logging.Error(ctx, "failed to list entries", "error", err, "filter", filter)
		return core.PaginatedResponse[Entry]{}, err
	}
	totalCount, err := s.repo.CountEntries(ctx, filter)
	if err != nil {
		logging.Error(ctx, "failed to count entries", "error", err, "filter", filter)
		return core.PaginatedResponse[Entry]{}, err
	}
	if entries == nil {
		entries = []Entry{}
	}
	logging.Info(ctx, "successfully listed entries", "user_id", userID, "count", len(entries), "limit", limit, "offset", offset, "deleted", deleted)
	return core.PaginatedResponse[Entry]{
		Items:      entries,
		TotalCount: totalCount,
	}, nil
}

func (s *Service) GetEntry(ctx context.Context, userID, id uint) (*Entry, error) {
	logging.Debug(ctx, "GetEntry service called", "user_id", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeAll)
	logging.Debug(ctx, "getting entry", "filter", filter)
	entry, err := s.repo.GetEntry(ctx, filter)
	if err != nil {
		logging.Error(ctx, "failed to get entry", "error", err, "filter", filter)
		return nil, err
	}
	logging.Info(ctx, "successfully retrieved entry", "user_id", userID, "id", id)
	return entry, nil
}

func (s *Service) CreateEntry(ctx context.Context, userID uint, req EditEntryRequest) (*Entry, error) {
	logging.Debug(ctx, "CreateEntry service called", "user_id", userID)
	entry, err := s.repo.SaveEntry(ctx, &Entry{
		UserID:  userID,
		Feeling: req.Feeling,
		Note:    req.Note,
	})
	if err != nil {
		logging.Error(ctx, "failed to create entry", "error", err, "user_id", userID)
		return nil, err
	}
	logging.Info(ctx, "successfully created entry", "user_id", userID, "id", entry.ID)
	return entry, nil
}

func (s *Service) UpdateEntry(ctx context.Context, userID, id uint, req EditEntryRequest) (*Entry, error) {
	logging.Debug(ctx, "UpdateEntry service called", "user_id", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id)
	logging.Debug(ctx, "updating entry", "filter", filter)
	entry, err := s.repo.GetEntry(ctx, filter)
	if err != nil {
		logging.Error(ctx, "failed to get entry for update", "error", err, "filter", filter)
		return nil, err
	}
	entry.Feeling = req.Feeling
	entry.Note = req.Note
	updatedEntry, err := s.repo.SaveEntry(ctx, entry)
	if err != nil {
		logging.Error(ctx, "failed to update entry", "error", err, "filter", filter)
		return nil, err
	}
	logging.Info(ctx, "successfully updated entry", "user_id", userID, "id", id)
	return updatedEntry, nil
}

func (s *Service) DeleteEntry(ctx context.Context, userID, id uint) (*Entry, error) {
	logging.Debug(ctx, "DeleteEntry service called", "user_id", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id)
	logging.Debug(ctx, "deleting entry", "filter", filter)
	entry, err := s.repo.DeleteEntry(ctx, filter)
	if err != nil {
		logging.Error(ctx, "failed to delete entry", "error", err, "filter", filter)
		return nil, err
	}
	logging.Info(ctx, "successfully deleted entry", "user_id", userID, "id", id)
	return entry, nil
}

func (s *Service) RestoreEntry(ctx context.Context, userID, id uint) (*Entry, error) {
	logging.Debug(ctx, "RestoreEntry service called", "user_id", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeDeletedOnly)
	logging.Debug(ctx, "restoring entry", "filter", filter)
	entry, err := s.repo.GetEntry(ctx, filter)
	if err != nil {
		logging.Error(ctx, "failed to get entry for restore", "error", err, "filter", filter)
		return nil, err
	}
	entry.DeletedAt.Valid = false
	restoredEntry, err := s.repo.SaveEntry(ctx, entry)
	if err != nil {
		logging.Error(ctx, "failed to restore entry", "error", err, "user_id", userID, "id", id)
		return nil, err
	}
	logging.Info(ctx, "successfully restored entry", "user_id", userID, "id", id)
	return restoredEntry, nil
}
