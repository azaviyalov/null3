package core

import "errors"

var (
	ErrInvalidItem  = errors.New("invalid item")
	ErrItemNotFound = errors.New("item not found")
)
