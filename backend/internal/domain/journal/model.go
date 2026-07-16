package journal

import (
	"time"

	"gorm.io/gorm"
)

type MoodEntry struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `json:"user_id"`
	Feeling      string         `json:"feeling" validate:"required"`
	Emoji        string         `json:"emoji,omitempty"`
	Note         string         `json:"note,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	DiaryEntries []DiaryEntry   `gorm:"many2many:mood_entry_diary_entries;->" json:"diary_entries,omitempty"`
}

type DiaryEntry struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `json:"user_id"`
	Title       string         `json:"title,omitempty"`
	Markdown    string         `json:"markdown" validate:"required"`
	OccurredAt  time.Time      `json:"occurred_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	MoodEntries []MoodEntry    `gorm:"many2many:mood_entry_diary_entries;" json:"mood_entries,omitempty"`
}
