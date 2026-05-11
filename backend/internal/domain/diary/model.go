package diary

import (
	"time"

	"gorm.io/gorm"
)

type Entry struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `json:"user_id"`
	Title      string         `json:"title,omitempty"`
	Markdown   string         `json:"markdown" validate:"required"`
	OccurredAt time.Time      `json:"occurred_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (e *Entry) TableName() string {
	return "diary_entries"
}

type MoodEntryLink struct {
	DiaryEntryID uint      `gorm:"primaryKey;autoIncrement:false" json:"diary_entry_id"`
	MoodEntryID  uint      `gorm:"primaryKey;autoIncrement:false" json:"mood_entry_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (m *MoodEntryLink) TableName() string {
	return "diary_entry_mood_entry_links"
}
