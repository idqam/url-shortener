package model

import (
	"time"
)

type URL struct {
	ID          string    `json:"id"`
	UserID      *string   `json:"user_id,omitempty"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	IsPublic    bool      `json:"is_public"`
	ClickCount  int       `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
}
