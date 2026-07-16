package journal

import (
	"time"

	"gorm.io/gorm"
)

type MoodEditRecordRequest struct {
	Feeling string `json:"feeling" validate:"required"`
	Emoji   string `json:"emoji,omitempty"`
	Note    string `json:"note,omitempty"`
}

type MoodRecordResponse struct {
	ID              uint                     `json:"id"`
	UserID          uint                     `json:"user_id"`
	Feeling         string                   `json:"feeling"`
	Emoji           string                   `json:"emoji,omitempty"`
	Note            string                   `json:"note,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	DeletedAt       gorm.DeletedAt           `json:"deleted_at"`
	DiaryEntryLinks []DiaryEntryLinkResponse `json:"diary_entry_links,omitempty"`
}

func NewMoodRecordResponse(entry *MoodRecord) MoodRecordResponse {
	return MoodRecordResponse{
		ID:              entry.ID,
		UserID:          entry.UserID,
		Feeling:         entry.Feeling,
		Emoji:           entry.Emoji,
		Note:            entry.Note,
		CreatedAt:       entry.CreatedAt,
		UpdatedAt:       entry.UpdatedAt,
		DeletedAt:       entry.DeletedAt,
		DiaryEntryLinks: NewDiaryEntryLinkResponses(entry.DiaryEntries),
	}
}

type DiaryEntryLinkResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title,omitempty"`
	Preview    string    `json:"preview,omitempty"`
	OccurredAt time.Time `json:"occurred_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewDiaryEntryLinkResponses(entries []DiaryEntry) []DiaryEntryLinkResponse {
	if len(entries) == 0 {
		return nil
	}

	result := make([]DiaryEntryLinkResponse, 0, len(entries))
	for _, entry := range entries {
		result = append(result, DiaryEntryLinkResponse{
			ID:         entry.ID,
			Title:      entry.Title,
			Preview:    MarkdownPreview(entry.Markdown),
			OccurredAt: entry.OccurredAt,
			CreatedAt:  entry.CreatedAt,
			UpdatedAt:  entry.UpdatedAt,
		})
	}

	return result
}

type DiaryEditEntryRequest struct {
	Title      string     `json:"title,omitempty"`
	Markdown   string     `json:"markdown" validate:"required"`
	OccurredAt *time.Time `json:"occurred_at" validate:"required"`
}

type DiaryEntryResponse struct {
	ID                    uint                 `json:"id"`
	UserID                uint                 `json:"user_id"`
	Title                 string               `json:"title,omitempty"`
	Markdown              string               `json:"markdown"`
	Preview               string               `json:"preview,omitempty"`
	OccurredAt            time.Time            `json:"occurred_at"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
	DeletedAt             gorm.DeletedAt       `json:"deleted_at"`
	ReferencedMoodRecords []MoodRecordResponse `json:"referenced_mood_records,omitempty"`
}

func NewDiaryEntryResponse(entry *DiaryEntry) DiaryEntryResponse {
	return DiaryEntryResponse{
		ID:                    entry.ID,
		UserID:                entry.UserID,
		Title:                 entry.Title,
		Markdown:              entry.Markdown,
		Preview:               MarkdownPreview(entry.Markdown),
		OccurredAt:            entry.OccurredAt,
		CreatedAt:             entry.CreatedAt,
		UpdatedAt:             entry.UpdatedAt,
		DeletedAt:             entry.DeletedAt,
		ReferencedMoodRecords: NewReferencedMoodRecordResponses(entry.MoodRecords),
	}
}

func NewReferencedMoodRecordResponses(entries []MoodRecord) []MoodRecordResponse {
	if len(entries) == 0 {
		return nil
	}

	result := make([]MoodRecordResponse, 0, len(entries))
	for _, entry := range entries {
		result = append(result, MoodRecordResponse{
			ID:        entry.ID,
			UserID:    entry.UserID,
			Feeling:   entry.Feeling,
			Emoji:     entry.Emoji,
			Note:      entry.Note,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
			DeletedAt: entry.DeletedAt,
		})
	}

	return result
}
