package mood

import (
	"github.com/azaviyalov/null3/backend/internal/core"
	"gorm.io/gorm"
)

type EntryFilter struct {
	ID          *uint
	UserID      *uint
	DeletedMode core.DeletedFilterMode
}

func NewEntryFilter() *EntryFilter {
	return &EntryFilter{
		ID:          nil,
		UserID:      nil,
		DeletedMode: core.DeletedModeNonDeleted,
	}
}

func (f *EntryFilter) WithID(id uint) *EntryFilter {
	f.ID = &id
	return f
}

func (f *EntryFilter) WithUserID(userID uint) *EntryFilter {
	f.UserID = &userID
	return f
}

func (f *EntryFilter) WithDeletedMode(mode core.DeletedFilterMode) *EntryFilter {
	f.DeletedMode = mode
	return f
}

func (f EntryFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		db = db.Where("user_id = ?", *f.UserID)
	}
	switch f.DeletedMode {
	case core.DeletedModeNonDeleted:
		// default GORM behavior: only non-deleted
		// do nothing
	case core.DeletedModeDeletedOnly:
		db = db.Unscoped().Where("deleted_at IS NOT NULL")
	case core.DeletedModeAll:
		db = db.Unscoped()
	default:
		// unknown mode, treat as non-deleted
		// default GORM behavior: only non-deleted
		// do nothing
	}
	return db
}
