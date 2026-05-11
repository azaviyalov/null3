package mood

import (
	"time"

	"gorm.io/gorm"
)

type EditEntryRequest struct {
	Feeling string `json:"feeling" validate:"required"`
	Emoji   string `json:"emoji,omitempty"`
	Note    string `json:"note,omitempty"`
}

type EntryResponse struct {
	ID              uint             `json:"id"`
	UserID          uint             `json:"user_id"`
	Feeling         string           `json:"feeling"`
	Emoji           string           `json:"emoji,omitempty"`
	Note            string           `json:"note,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `json:"deleted_at"`
	DiaryEntryLinks []DiaryEntryLink `json:"diary_entry_links"`
}

func NewEntryResponse(entry *Entry, diaryEntryLinks []DiaryEntryLink) EntryResponse {
	if diaryEntryLinks == nil {
		diaryEntryLinks = []DiaryEntryLink{}
	}

	return EntryResponse{
		ID:              entry.ID,
		UserID:          entry.UserID,
		Feeling:         entry.Feeling,
		Emoji:           entry.Emoji,
		Note:            entry.Note,
		CreatedAt:       entry.CreatedAt,
		UpdatedAt:       entry.UpdatedAt,
		DeletedAt:       entry.DeletedAt,
		DiaryEntryLinks: diaryEntryLinks,
	}
}
