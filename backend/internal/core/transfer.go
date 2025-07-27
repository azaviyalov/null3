package core

type PaginatedResponse[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
}
