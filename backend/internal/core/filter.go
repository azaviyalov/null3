package core

type DeletedFilterMode int

const (
	DeletedModeNonDeleted  DeletedFilterMode = iota // default: only non-deleted
	DeletedModeDeletedOnly                          // only deleted
	DeletedModeAll                                  // both deleted and non-deleted
)
