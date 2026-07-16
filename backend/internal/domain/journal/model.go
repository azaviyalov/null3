package journal

import (
	"time"

	"gorm.io/gorm"
)

type MoodRecord struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `json:"user_id"`
	Feeling      string         `json:"feeling" validate:"required"`
	Emoji        string         `json:"emoji,omitempty"`
	Note         string         `json:"note,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	DiaryEntries []DiaryEntry   `gorm:"many2many:mood_record_diary_entries;joinForeignKey:MoodRecordID;joinReferences:DiaryEntryID;->" json:"diary_entries,omitempty"`
}

func (MoodRecord) TableName() string {
	return "mood_records"
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
	MoodRecords []MoodRecord   `gorm:"many2many:mood_record_diary_entries;joinForeignKey:DiaryEntryID;joinReferences:MoodRecordID" json:"mood_records,omitempty"`
}
