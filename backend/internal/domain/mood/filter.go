package mood

import (
	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
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

func (f EntryFilter) ToFieldValue() logging.FieldValue {
	idVal := logging.NewStringValue("null")
	if f.ID != nil {
		idVal = logging.NewUint64Value(uint64(*f.ID))
	}
	userIDVal := logging.NewStringValue("null")
	if f.UserID != nil {
		userIDVal = logging.NewUint64Value(uint64(*f.UserID))
	}
	return logging.CombineFields(
		logging.NewField("id", idVal),
		logging.NewField("user_id", userIDVal),
		logging.NewField("deleted_mode", logging.NewStringValue(f.DeletedMode.String())),
	)
}
