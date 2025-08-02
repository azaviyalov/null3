package core

type DeletedFilterMode int

const (
	DeletedModeNonDeleted  DeletedFilterMode = iota // default: only non-deleted
	DeletedModeDeletedOnly                          // only deleted
	DeletedModeAll                                  // both deleted and non-deleted
)

func (m DeletedFilterMode) String() string {
	switch m {
	case DeletedModeNonDeleted:
		return "NonDeleted"
	case DeletedModeDeletedOnly:
		return "DeletedOnly"
	case DeletedModeAll:
		return "All"
	default:
		return "Unknown"
	}
}
