package journal

import (
	"github.com/azaviyalov/null3/backend/internal/core"
	"gorm.io/gorm"
)

type MoodEntryFilter struct {
	ID          *uint
	UserID      *uint
	DeletedMode core.DeletedFilterMode
}

func NewMoodEntryFilter() *MoodEntryFilter {
	return &MoodEntryFilter{DeletedMode: core.DeletedModeNonDeleted}
}

func (f *MoodEntryFilter) WithID(id uint) *MoodEntryFilter {
	f.ID = &id
	return f
}

func (f *MoodEntryFilter) WithUserID(userID uint) *MoodEntryFilter {
	f.UserID = &userID
	return f
}

func (f *MoodEntryFilter) WithDeletedMode(mode core.DeletedFilterMode) *MoodEntryFilter {
	f.DeletedMode = mode
	return f
}

func (f MoodEntryFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		db = db.Where("user_id = ?", *f.UserID)
	}
	switch f.DeletedMode {
	case core.DeletedModeNonDeleted:
	case core.DeletedModeDeletedOnly:
		db = db.Unscoped().Where("deleted_at IS NOT NULL")
	case core.DeletedModeAll:
		db = db.Unscoped()
	}
	return db
}

type DiaryEntryFilter struct {
	ID          *uint
	UserID      *uint
	DeletedMode core.DeletedFilterMode
}

func NewDiaryEntryFilter() *DiaryEntryFilter {
	return &DiaryEntryFilter{DeletedMode: core.DeletedModeNonDeleted}
}

func (f *DiaryEntryFilter) WithID(id uint) *DiaryEntryFilter {
	f.ID = &id
	return f
}

func (f *DiaryEntryFilter) WithUserID(userID uint) *DiaryEntryFilter {
	f.UserID = &userID
	return f
}

func (f *DiaryEntryFilter) WithDeletedMode(mode core.DeletedFilterMode) *DiaryEntryFilter {
	f.DeletedMode = mode
	return f
}

func (f DiaryEntryFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		db = db.Where("user_id = ?", *f.UserID)
	}
	switch f.DeletedMode {
	case core.DeletedModeNonDeleted:
	case core.DeletedModeDeletedOnly:
		db = db.Unscoped().Where("deleted_at IS NOT NULL")
	case core.DeletedModeAll:
		db = db.Unscoped()
	}
	return db
}
