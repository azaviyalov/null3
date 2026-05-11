package mood

import "time"

type DiaryEntryLink struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title,omitempty"`
	Preview    string    `json:"preview,omitempty"`
	OccurredAt time.Time `json:"occurred_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
