package journal

import (
	"github.com/azaviyalov/null3/backend/internal/core"
	"gorm.io/gorm"
)

type MoodRecordFilter struct {
	ID          *uint
	UserID      *uint
	DeletedMode core.DeletedFilterMode
}

func NewMoodRecordFilter() *MoodRecordFilter {
	return &MoodRecordFilter{DeletedMode: core.DeletedModeNonDeleted}
}

func (f *MoodRecordFilter) WithID(id uint) *MoodRecordFilter {
	f.ID = &id
	return f
}

func (f *MoodRecordFilter) WithUserID(userID uint) *MoodRecordFilter {
	f.UserID = &userID
	return f
}

func (f *MoodRecordFilter) WithDeletedMode(mode core.DeletedFilterMode) *MoodRecordFilter {
	f.DeletedMode = mode
	return f
}

func (f MoodRecordFilter) Apply(db *gorm.DB) *gorm.DB {
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
