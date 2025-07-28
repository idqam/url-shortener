package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"url-shortener-go-backend/internal/model"
)

type URLRepositoryImpl struct {
	*SupabaseRepository
}

func NewURLRepository(baseRepo *SupabaseRepository) URLRepository {
	return &URLRepositoryImpl{baseRepo}
}

func (u *URLRepositoryImpl) GetURLByShortCode(shortcode string) (*model.URL, error) {
	resp, _, err := u.Client.
		From("urls").
		Select("*", "exact", false).
		Eq("short_code", shortcode).
		Execute()

	if err != nil {
		log.Printf("[GetURLByShortCode] Failed to fetch: %v", err)
		return nil, fmt.Errorf("failed to fetch URL by shortcode: %w", err)
	}
	if len(resp) == 0 || string(resp) == "[]" {
		log.Printf("[GetURLByShortCode] No result found for shortcode: %s", shortcode)
		return nil, nil
	}

	var results []model.URL
	if err := json.Unmarshal(resp, &results); err != nil {
		log.Printf("[GetURLByShortCode] Unmarshal error: %v", err)
		return nil, fmt.Errorf("failed to unmarshal URL data: %w", err)
	}

	return &results[0], nil
}

func (u *URLRepositoryImpl) GetURLByOriginalURL(original string) (*model.URL, error) {
	resp, _, err := u.Client.
		From("urls").
		Select("*", "exact", false).
		Eq("original_url", original).
		Execute()

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") ||
			strings.Contains(err.Error(), "No rows") {
			log.Printf("[GetURLByOriginalURL] No match for: %s", original)
			return nil, nil
		}
		log.Printf("[GetURLByOriginalURL] Query error: %v", err)
		return nil, fmt.Errorf("failed to fetch URL by original: %w", err)
	}
	if len(resp) == 0 || string(resp) == "[]" {
		log.Printf("[GetURLByOriginalURL] Empty response for: %s", original)
		return nil, nil
	}

	var results []model.URL
	if err := json.Unmarshal(resp, &results); err != nil {
		log.Printf("[GetURLByOriginalURL] Unmarshal error: %v", err)
		return nil, fmt.Errorf("failed to unmarshal URL data: %w", err)
	}
	return &results[0], nil
}

func (u *URLRepositoryImpl) GetUserUrls(user model.User) ([]model.URL, error) {
	resp, _, err := u.Client.
		From("urls").
		Select("*", "exact", false).
		Eq("user_id", user.ID).
		Execute()

	if err != nil {
		if strings.Contains(err.Error(), "No rows") ||
			strings.Contains(strings.ToLower(err.Error()), "not found") {
			log.Printf("[GetUserUrls] No URLs found for user ID: %s", user.ID)
			return []model.URL{}, nil
		}
		log.Printf("[GetUserUrls] Query error: %v", err)
		return nil, fmt.Errorf("failed to get user urls: %w", err)
	}

	var urls []model.URL
	if err := json.Unmarshal(resp, &urls); err != nil {
		log.Printf("[GetUserUrls] Unmarshal error: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return urls, nil
}

func (u *URLRepositoryImpl) IncrementClickCount(shortcode string) error {
	err := u.Client.Rpc("increment_click_count", "", map[string]any{
		"sc": shortcode,
	})

	if err != "" {
		log.Printf("[IncrementClickCount] RPC failed for shortcode=%s: %s", shortcode, err)
		return fmt.Errorf("failed to increment click amount: %s", err)
	}
	return nil
}

func (u *URLRepositoryImpl) SaveURL(url *model.URL) error {
	data := map[string]interface{}{
		"original_url": url.OriginalURL,
		"short_code":   url.ShortCode,
		"is_public":    url.IsPublic,
		"click_count":  url.ClickCount,
	}

	if url.UserID != nil && *url.UserID != "" {
		data["user_id"] = *url.UserID
	}

	resp, _, err := u.Client.
		From("urls").
		Insert(data, false, "", "", "").
		Execute()

	if err != nil {
		log.Printf("[SaveURL] Insert failed: %v", err)
		return fmt.Errorf("supabase insert failed: %w", err)
	}

	if len(resp) == 0 {
		log.Printf("[SaveURL] Empty insert response")
		return fmt.Errorf("failed to save URL (empty response)")
	}

	var inserted []model.URL
	if err := json.Unmarshal(resp, &inserted); err != nil || len(inserted) == 0 {
		log.Printf("[SaveURL] Failed to decode insert response: %v", err)
		var errResp struct {
			Message string `json:"message"`
			Code    string `json:"code"`
			Details string `json:"details"`
			Hint    string `json:"hint"`
		}
		if json.Unmarshal(resp, &errResp) == nil && errResp.Message != "" {
			msg := strings.ToLower(errResp.Message)
			if strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint") {
				log.Printf("[SaveURL] Unique constraint violation: %s", errResp.Message)
				return ErrUniqueViolation
			}
			return fmt.Errorf("supabase error: %s", errResp.Message)
		}
		return fmt.Errorf("failed to decode inserted URL: %s", string(resp))
	}

	*url = inserted[0]
	url.PopulateShortURL()
	log.Printf("[SaveURL] Insert successful: id=%s, short_code=%s", url.ID, url.ShortCode)
	return nil
}
