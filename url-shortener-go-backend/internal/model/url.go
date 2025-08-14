package model

import (
	"fmt"
	"os"
	"strings"
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
	ShortURL    string    `json:"short_url"`
}

type URLSubset struct {
	Original_URL string `json:"original_url"`
	Short_Code   string `json:"short_code"`
	Is_Public    bool   `json:"is_public"`
	Click_Count  int    `json:"click_count"`
}

func (u *URL) PopulateShortURL() {
	baseDomain := os.Getenv("SHORT_DOMAIN")
	if baseDomain == "" {
		baseDomain = "http://localhost:8080"
	}
	u.ShortURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(baseDomain, "/"), u.ShortCode)
}
