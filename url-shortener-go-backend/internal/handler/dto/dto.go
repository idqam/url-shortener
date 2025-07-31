package dto

import "time"

type ShortenRequest struct {
	URL        string  `json:"url"`
	IsPublic   bool    `json:"is_public"`
	UserID     *string `json:"user_id,omitempty"`
	CodeLength int8    `json:"code_length"`
}

type GetUrlsRequest struct {
	UserID string `json:"user_id"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type GetUserByEmailResponse struct {
	ID string `json:"id"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`

}

type ShortenResponse struct {
	ID          string `json:"id"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type URLResponse struct {
	ID          string  `json:"id"`
	UserID      *string `json:"user_id,omitempty"`
	OriginalURL string  `json:"original_url"`
	ShortCode   string  `json:"short_code"`
	ShortURL    string  `json:"short_url"`
	IsPublic    bool    `json:"is_public"`
	ClickCount  int     `json:"click_count"`
	CreatedAt   string  `json:"created_at"`
}

type GetUrlsResponse struct {
	URLs []URLResponse `json:"urls"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
