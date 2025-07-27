package mood

import (
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

func (s *Service) ListEntries(userID uint, limit, offset int) (core.PaginatedResponse[Entry], error) {
	filter := NewEntryFilter().WithUserID(userID)

	entries, err := s.repo.ListEntries(filter, limit, offset)
	if err != nil {
		return core.PaginatedResponse[Entry]{}, nil
	}

	totalCount, err := s.repo.CountEntries(filter)
	if err != nil {
		return core.PaginatedResponse[Entry]{}, err
	}

	if entries == nil {
		entries = []Entry{}
	}

	return core.PaginatedResponse[Entry]{
		Items:      entries,
		TotalCount: totalCount,
	}, nil
}

func (s *Service) GetEntry(userID, id uint) (*Entry, error) {
	return s.repo.GetEntry(NewEntryFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeAll))
}

func (s *Service) CreateEntry(userID uint, req EditEntryRequest) (*Entry, error) {
	return s.repo.SaveEntry(&Entry{
		UserID:  userID,
		Feeling: req.Feeling,
		Note:    req.Note,
	})
}

func (s *Service) UpdateEntry(userID, id uint, req EditEntryRequest) (*Entry, error) {
	entry, err := s.repo.GetEntry(NewEntryFilter().WithUserID(userID).WithID(id))
	if err != nil {
		return nil, err
	}

	entry.Feeling = req.Feeling
	entry.Note = req.Note

	return s.repo.SaveEntry(entry)
}

func (s *Service) DeleteEntry(userID, id uint) (*Entry, error) {
	return s.repo.DeleteEntry(NewEntryFilter().WithUserID(userID).WithID(id))
}

func (s *Service) RestoreEntry(userID, id uint) (*Entry, error) {
	entry, err := s.repo.GetEntry(NewEntryFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeDeletedOnly))
	if err != nil {
		return nil, err
	}

	entry.DeletedAt.Valid = false

	return s.repo.SaveEntry(entry)
}
