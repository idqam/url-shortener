package dto


type CreateURLRequest struct {
	URL        string  `json:"url"`
	IsPublic   bool    `json:"is_public"`
	UserID     *string `json:"user_id,omitempty"`
	CodeLength int8    `json:"code_length"`
}

type GetUrlsRequest struct {
	UserID string `json:"user_id"`
}


type CreateURLResponse struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
}

type URLResponse struct {
	ID          string  `json:"id"`
	UserID      *string `json:"user_id,omitempty"`
	OriginalURL string  `json:"original_url"`
	ShortCode   string  `json:"short_code"`
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

type ShortenRequest struct {
    URL        string  `json:"url"`
    IsPublic   bool    `json:"is_public"`
    UserID     *string `json:"user_id,omitempty"`
    CodeLength int8    `json:"code_length"`
}

type ShortenResponse struct {
    ID          string `json:"id"`
    ShortCode   string `json:"short_code"`
    OriginalURL string `json:"original_url"`
    ShortURL    string `json:"short_url"`
}