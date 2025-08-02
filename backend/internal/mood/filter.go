package mood

import (
	"log/slog"

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

func (f EntryFilter) LogValue() slog.Value {
	var idVal, userIDVal slog.Value
	if f.ID != nil {
		idVal = slog.Uint64Value(uint64(*f.ID))
	} else {
		idVal = slog.StringValue("<nil>")
	}
	if f.UserID != nil {
		userIDVal = slog.Uint64Value(uint64(*f.UserID))
	} else {
		userIDVal = slog.StringValue("<nil>")
	}
	return slog.GroupValue(
		slog.Attr{Key: "ID", Value: idVal},
		slog.Attr{Key: "UserID", Value: userIDVal},
		slog.String("DeletedMode", f.DeletedMode.String()),
	)
}
