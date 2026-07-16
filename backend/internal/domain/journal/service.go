package journal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListMoodRecords(ctx context.Context, userID uint, limit, offset int, deleted bool) (core.Page[MoodRecord], error) {
	filter := NewMoodRecordFilter().WithUserID(userID)
	if deleted {
		filter = filter.WithDeletedMode(core.DeletedModeDeletedOnly)
	}

	entries, err := s.repo.ListMoodRecords(ctx, filter, limit, offset)
	if err != nil {
		return core.Page[MoodRecord]{}, err
	}
	totalCount, err := s.repo.CountMoodRecords(ctx, filter)
	if err != nil {
		return core.Page[MoodRecord]{}, err
	}
	if entries == nil {
		entries = []MoodRecord{}
	}
	return core.Page[MoodRecord]{Items: entries, TotalCount: totalCount}, nil
}

func (s *Service) GetMoodRecord(ctx context.Context, userID, id uint) (*MoodRecord, error) {
	filter := NewMoodRecordFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeAll)
	return s.repo.GetMoodRecord(ctx, filter)
}

func (s *Service) CreateMoodRecord(ctx context.Context, userID uint, req MoodEditRecordRequest) (*MoodRecord, error) {
	return s.repo.SaveMoodRecord(ctx, &MoodRecord{
		UserID:  userID,
		Feeling: req.Feeling,
		Emoji:   req.Emoji,
		Note:    req.Note,
	})
}

func (s *Service) UpdateMoodRecord(ctx context.Context, userID, id uint, req MoodEditRecordRequest) (*MoodRecord, error) {
	filter := NewMoodRecordFilter().WithUserID(userID).WithID(id)
	entry, err := s.repo.GetMoodRecord(ctx, filter)
	if err != nil {
		return nil, err
	}
	entry.Feeling = req.Feeling
	entry.Emoji = req.Emoji
	entry.Note = req.Note
	return s.repo.SaveMoodRecord(ctx, entry)
}

func (s *Service) DeleteMoodRecord(ctx context.Context, userID, id uint) (*MoodRecord, error) {
	filter := NewMoodRecordFilter().WithUserID(userID).WithID(id)
	return s.repo.DeleteMoodRecord(ctx, filter)
}

func (s *Service) RestoreMoodRecord(ctx context.Context, userID, id uint) (*MoodRecord, error) {
	filter := NewMoodRecordFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeDeletedOnly)
	entry, err := s.repo.GetMoodRecord(ctx, filter)
	if err != nil {
		return nil, err
	}
	entry.DeletedAt.Valid = false
	return s.repo.SaveMoodRecord(ctx, entry)
}

func (s *Service) ListDiaryEntries(ctx context.Context, userID uint, limit, offset int, deleted bool) (core.Page[DiaryEntry], error) {
	filter := NewDiaryEntryFilter().WithUserID(userID)
	if deleted {
		filter = filter.WithDeletedMode(core.DeletedModeDeletedOnly)
	}

	entries, err := s.repo.ListDiaryEntries(ctx, filter, limit, offset)
	if err != nil {
		return core.Page[DiaryEntry]{}, err
	}
	totalCount, err := s.repo.CountDiaryEntries(ctx, filter)
	if err != nil {
		return core.Page[DiaryEntry]{}, err
	}
	if entries == nil {
		entries = []DiaryEntry{}
	}
	return core.Page[DiaryEntry]{Items: entries, TotalCount: totalCount}, nil
}

func (s *Service) GetDiaryEntry(ctx context.Context, userID, id uint) (*DiaryEntry, error) {
	filter := NewDiaryEntryFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeAll)
	return s.repo.GetDiaryEntry(ctx, filter)
}

func (s *Service) CreateDiaryEntry(ctx context.Context, userID uint, req DiaryEditEntryRequest) (*DiaryEntry, error) {
	title, markdown, occurredAt, err := normalizeDiaryRequest(req)
	if err != nil {
		return nil, err
	}

	moodRecords, err := s.resolveDiaryMoodRecords(ctx, userID, markdown)
	if err != nil {
		return nil, err
	}

	return s.repo.SaveDiaryEntry(ctx, &DiaryEntry{
		UserID:      userID,
		Title:       title,
		Markdown:    markdown,
		OccurredAt:  occurredAt,
		MoodRecords: moodRecords,
	})
}

func (s *Service) UpdateDiaryEntry(ctx context.Context, userID, id uint, req DiaryEditEntryRequest) (*DiaryEntry, error) {
	filter := NewDiaryEntryFilter().WithUserID(userID).WithID(id)
	entry, err := s.repo.GetDiaryEntry(ctx, filter)
	if err != nil {
		return nil, err
	}
	title, markdown, occurredAt, err := normalizeDiaryRequest(req)
	if err != nil {
		return nil, err
	}

	moodRecords, err := s.resolveDiaryMoodRecords(ctx, userID, markdown)
	if err != nil {
		return nil, err
	}

	entry.Title = title
	entry.Markdown = markdown
	entry.OccurredAt = occurredAt
	entry.MoodRecords = moodRecords
	return s.repo.SaveDiaryEntry(ctx, entry)
}

func (s *Service) DeleteDiaryEntry(ctx context.Context, userID, id uint) (*DiaryEntry, error) {
	filter := NewDiaryEntryFilter().WithUserID(userID).WithID(id)
	return s.repo.DeleteDiaryEntry(ctx, filter)
}

func (s *Service) RestoreDiaryEntry(ctx context.Context, userID, id uint) (*DiaryEntry, error) {
	filter := NewDiaryEntryFilter().WithUserID(userID).WithID(id).WithDeletedMode(core.DeletedModeDeletedOnly)
	entry, err := s.repo.GetDiaryEntry(ctx, filter)
	if err != nil {
		return nil, err
	}
	entry.DeletedAt.Valid = false
	return s.repo.SaveDiaryEntry(ctx, entry)
}

func normalizeDiaryRequest(req DiaryEditEntryRequest) (string, string, time.Time, error) {
	title := strings.TrimSpace(req.Title)
	markdown := strings.TrimSpace(req.Markdown)
	if markdown == "" {
		return "", "", time.Time{}, fmt.Errorf("%w: markdown is required", core.ErrInvalidItem)
	}

	if req.OccurredAt == nil || req.OccurredAt.IsZero() {
		return "", "", time.Time{}, fmt.Errorf("%w: occurred_at is required", core.ErrInvalidItem)
	}

	occurredAt := req.OccurredAt.UTC()
	if occurredAt.After(time.Now().UTC()) {
		return "", "", time.Time{}, fmt.Errorf("%w: occurred_at cannot be in the future", core.ErrInvalidItem)
	}

	return title, markdown, occurredAt, nil
}

func (s *Service) resolveDiaryMoodRecords(ctx context.Context, userID uint, markdown string) ([]MoodRecord, error) {
	ids, err := ExtractMoodRecordIDs(markdown)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid mood references", core.ErrInvalidItem)
	}

	moodRecords, err := s.repo.ListMoodRecordsByIDs(ctx, userID, ids)
	if err != nil {
		return nil, err
	}

	if len(moodRecords) != len(ids) {
		return nil, fmt.Errorf("%w: one or more referenced mood records do not exist", core.ErrInvalidItem)
	}

	return moodRecords, nil
}
