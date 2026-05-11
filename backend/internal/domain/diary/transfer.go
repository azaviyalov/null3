package diary

import (
	"time"

	mooddomain "github.com/azaviyalov/null3/backend/internal/domain/mood"
	"gorm.io/gorm"
)

type EditEntryRequest struct {
	Title      string     `json:"title,omitempty"`
	Markdown   string     `json:"markdown" validate:"required"`
	OccurredAt *time.Time `json:"occurred_at" validate:"required"`
}

type ReferencedMoodEntry struct {
	ID        uint           `json:"id"`
	Feeling   string         `json:"feeling"`
	Emoji     string         `json:"emoji,omitempty"`
	Note      string         `json:"note,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func NewReferencedMoodEntry(entry mooddomain.Entry) ReferencedMoodEntry {
	return ReferencedMoodEntry{
		ID:        entry.ID,
		Feeling:   entry.Feeling,
		Emoji:     entry.Emoji,
		Note:      entry.Note,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
		DeletedAt: entry.DeletedAt,
	}
}

type EntryResponse struct {
	ID                    uint                  `json:"id"`
	UserID                uint                  `json:"user_id"`
	Title                 string                `json:"title,omitempty"`
	Markdown              string                `json:"markdown"`
	Preview               string                `json:"preview,omitempty"`
	OccurredAt            time.Time             `json:"occurred_at"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
	DeletedAt             gorm.DeletedAt        `json:"deleted_at"`
	ReferencedMoodEntries []ReferencedMoodEntry `json:"referenced_mood_entries"`
}

func NewEntryResponse(entry *Entry, referencedMoodEntries []ReferencedMoodEntry) EntryResponse {
	if referencedMoodEntries == nil {
		referencedMoodEntries = []ReferencedMoodEntry{}
	}

	return EntryResponse{
		ID:                    entry.ID,
		UserID:                entry.UserID,
		Title:                 entry.Title,
		Markdown:              entry.Markdown,
		Preview:               MarkdownPreview(entry.Markdown),
		OccurredAt:            entry.OccurredAt,
		CreatedAt:             entry.CreatedAt,
		UpdatedAt:             entry.UpdatedAt,
		DeletedAt:             entry.DeletedAt,
		ReferencedMoodEntries: referencedMoodEntries,
	}
}
