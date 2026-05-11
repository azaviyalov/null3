package diary

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	mooddomain "github.com/azaviyalov/null3/backend/internal/domain/mood"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ListEntries(ctx context.Context, userID uint, limit, offset int, deleted bool) (core.Page[Entry], error) {
	logging.Debug(ctx, "ListEntries service called", "user_id", userID, "limit", limit, "offset", offset, "deleted", deleted)
	filter := NewEntryFilter().WithUserID(userID)
	if deleted {
		filter = filter.WithDeletedMode(core.DeletedModeDeletedOnly)
	}

	entries, err := s.repo.ListEntries(ctx, filter, limit, offset)
	if err != nil {
		logging.Error(ctx, "failed to list diary entries", "error", err, "filter", filter)
		return core.Page[Entry]{}, err
	}

	totalCount, err := s.repo.CountEntries(ctx, filter)
	if err != nil {
		logging.Error(ctx, "failed to count diary entries", "error", err, "filter", filter)
		return core.Page[Entry]{}, err
	}

	if entries == nil {
		entries = []Entry{}
	}

	logging.Info(ctx, "successfully listed diary entries", "user_id", userID, "count", len(entries), "limit", limit, "offset", offset, "deleted", deleted)
	return core.Page[Entry]{
		Items:      entries,
		TotalCount: totalCount,
	}, nil
}

func (s *Service) GetEntry(ctx context.Context, userID, id uint) (*Entry, error) {
	logging.Debug(ctx, "GetEntry service called", "user_id", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeAll)
	entry, err := s.repo.GetEntry(ctx, filter)
	if err != nil {
		logging.Error(ctx, "failed to get diary entry", "error", err, "filter", filter)
		return nil, err
	}

	logging.Info(ctx, "successfully retrieved diary entry", "user_id", userID, "id", id)
	return entry, nil
}

func (s *Service) CreateEntry(ctx context.Context, userID uint, req EditEntryRequest) (*Entry, error) {
	logging.Debug(ctx, "CreateEntry service called", "user_id", userID)

	title, markdown, occurredAt, referencedMoodEntryIDs, err := normalizeRequest(req)
	if err != nil {
		logging.Info(ctx, "invalid diary entry data", "error", err, "user_id", userID)
		return nil, err
	}

	entry, err := s.repo.SaveEntry(ctx, &Entry{
		UserID:     userID,
		Title:      title,
		Markdown:   markdown,
		OccurredAt: occurredAt,
	}, referencedMoodEntryIDs)
	if err != nil {
		logging.Error(ctx, "failed to create diary entry", "error", err, "user_id", userID)
		return nil, err
	}

	logging.Info(ctx, "successfully created diary entry", "user_id", userID, "id", entry.ID)
	return entry, nil
}

func (s *Service) UpdateEntry(ctx context.Context, userID, id uint, req EditEntryRequest) (*Entry, error) {
	logging.Debug(ctx, "UpdateEntry service called", "user_id", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id)
	entry, err := s.repo.GetEntry(ctx, filter)
	if err != nil {
		logging.Error(ctx, "failed to get diary entry for update", "error", err, "filter", filter)
		return nil, err
	}

	title, markdown, occurredAt, referencedMoodEntryIDs, err := normalizeRequest(req)
	if err != nil {
		logging.Info(ctx, "invalid diary entry update", "error", err, "user_id", userID, "id", id)
		return nil, err
	}

	entry.Title = title
	entry.Markdown = markdown
	entry.OccurredAt = occurredAt

	updatedEntry, err := s.repo.SaveEntry(ctx, entry, referencedMoodEntryIDs)
	if err != nil {
		logging.Error(ctx, "failed to update diary entry", "error", err, "filter", filter)
		return nil, err
	}

	logging.Info(ctx, "successfully updated diary entry", "user_id", userID, "id", id)
	return updatedEntry, nil
}

func (s *Service) DeleteEntry(ctx context.Context, userID, id uint) (*Entry, error) {
	logging.Debug(ctx, "DeleteEntry service called", "user_id", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id)
	entry, err := s.repo.DeleteEntry(ctx, filter)
	if err != nil {
		logging.Error(ctx, "failed to delete diary entry", "error", err, "filter", filter)
		return nil, err
	}

	logging.Info(ctx, "successfully deleted diary entry", "user_id", userID, "id", id)
	return entry, nil
}

func (s *Service) RestoreEntry(ctx context.Context, userID, id uint) (*Entry, error) {
	logging.Debug(ctx, "RestoreEntry service called", "user_id", userID, "id", id)
	filter := NewEntryFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeDeletedOnly)
	entry, err := s.repo.GetEntry(ctx, filter)
	if err != nil {
		logging.Error(ctx, "failed to get diary entry for restore", "error", err, "filter", filter)
		return nil, err
	}

	entry.DeletedAt.Valid = false
	referencedMoodEntryIDs, err := ExtractMoodEntryIDs(entry.Markdown)
	if err != nil {
		logging.Error(ctx, "failed to parse diary entry links for restore", "error", err, "user_id", userID, "id", id)
		return nil, fmt.Errorf("parse diary links: %w", err)
	}

	restoredEntry, err := s.repo.SaveEntry(ctx, entry, referencedMoodEntryIDs)
	if err != nil {
		logging.Error(ctx, "failed to restore diary entry", "error", err, "user_id", userID, "id", id)
		return nil, err
	}

	logging.Info(ctx, "successfully restored diary entry", "user_id", userID, "id", id)
	return restoredEntry, nil
}

func (s *Service) ListReferencedMoodEntries(ctx context.Context, userID, diaryEntryID uint) ([]ReferencedMoodEntry, error) {
	return s.repo.ListReferencedMoodEntries(ctx, userID, diaryEntryID)
}

func (s *Service) ListDiaryEntryLinks(ctx context.Context, userID, moodEntryID uint) ([]mooddomain.DiaryEntryLink, error) {
	return s.repo.ListDiaryEntryLinks(ctx, userID, moodEntryID)
}

func normalizeRequest(req EditEntryRequest) (string, string, time.Time, []uint, error) {
	title := strings.TrimSpace(req.Title)
	markdown := strings.TrimSpace(req.Markdown)
	if markdown == "" {
		return "", "", time.Time{}, nil, fmt.Errorf("%w: markdown is required", core.ErrInvalidItem)
	}

	if req.OccurredAt == nil || req.OccurredAt.IsZero() {
		return "", "", time.Time{}, nil, fmt.Errorf("%w: occurred_at is required", core.ErrInvalidItem)
	}

	occurredAt := req.OccurredAt.UTC()
	if occurredAt.After(time.Now().UTC()) {
		return "", "", time.Time{}, nil, fmt.Errorf("%w: occurred_at cannot be in the future", core.ErrInvalidItem)
	}

	referencedMoodEntryIDs, err := ExtractMoodEntryIDs(markdown)
	if err != nil {
		return "", "", time.Time{}, nil, fmt.Errorf("%w: invalid referenced mood entries", core.ErrInvalidItem)
	}

	return title, markdown, occurredAt, referencedMoodEntryIDs, nil
}
