package dto

type ShortenURLRequest struct {
	OriginalURL string `json:"url" validate:"required,url"`
	IsPublic    bool   `json:"is_public"`
	CodeLength  int8   `json:"code_length"`
}

type ShortenURLResponse struct {
	ID         string `json:"id"`
	ShortCode  string `json:"short_code"`
	ShortURL   string `json:"short_url"`
	CreatedAt  string `json:"created_at"`
	IsPublic   bool   `json:"is_public"`
	ClickCount int    `json:"click_count"`
}

type GetUserURLsResponse struct {
	URLs []ShortenURLResponse `json:"urls"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
	Field string `json:"field,omitempty"`
}

type GetURLByShortCodeResponse struct {
	OriginalURL string `json:"original_url"`
	ClickCount  int    `json:"click_count"`
}
