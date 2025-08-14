package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/utils"
)

type URLRepositoryImpl struct {
	*SupabaseRepository
}

func NewURLRepository(baseRepo *SupabaseRepository) URLRepository {
	return &URLRepositoryImpl{baseRepo}
}

func (u *URLRepositoryImpl) GetURLByShortCode(ctx context.Context, shortcode string) (*model.URL, error) {
	resp, _, err := u.Client.
		From("urls").
		Select("id, original_url, short_code, click_count, is_public, created_at", "exact", false).
		Eq("short_code", shortcode).
		Single().
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}

	var url model.URL
	if err := json.Unmarshal(resp, &url); err != nil {
		return nil, fmt.Errorf("failed to decode URL response: %w", err)
	}

	url.PopulateShortURL()

	return &url, nil
}

func (u *URLRepositoryImpl) GetUserUrls(ctx context.Context, userID string) ([]model.URL, error) {
	resp, _, err := u.Client.
		From("urls").
		Select("*", "exact", false).
		Eq("user_id", userID).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user URLs: %w", err)
	}

	var urls []model.URL
	if err := json.Unmarshal(resp, &urls); err != nil {
		return nil, fmt.Errorf("failed to decode URLs: %w", err)
	}

	for i := range urls {
		urls[i].PopulateShortURL()
	}

	return urls, nil
}

func (u *URLRepositoryImpl) IncrementClickCount(ctx context.Context, shortcode string) error {
	err := u.Client.Rpc("increment_click_count", "", map[string]any{
		"sc": shortcode,
	})

	if err != "" {
		log.Printf("[IncrementClickCount] RPC failed for shortcode=%s: %s", shortcode, err)
		return fmt.Errorf("failed to increment click count: %s", err)
	}

	return nil
}

func (u *URLRepositoryImpl) SaveURL(ctx context.Context, url *model.URL) error {
	data := map[string]interface{}{
		"original_url": url.OriginalURL,
		"short_code":   url.ShortCode,
		"is_public":    url.IsPublic,
		"click_count":  url.ClickCount,
	}

	if userID := url.UserID; userID != nil && *userID != "" {
		data["user_id"] = *userID
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
	if err := json.Unmarshal(resp, &inserted); err != nil {

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

	if len(inserted) == 0 {
		log.Printf("[SaveURL] Empty insert response")
		return fmt.Errorf("no URL returned after insert")
	}

	*url = inserted[0]

	url.PopulateShortURL()
	log.Printf("[SaveURL] Insert successful: id=%s, short_code=%s", url.ID, url.ShortCode)
	return nil
}

func (u *URLRepositoryImpl) SaveAnalytics(ctx context.Context, userID, urlID, referrer, deviceType string) error {
	data := map[string]interface{}{
		"url_id":      urlID,
		"referrer":    referrer,
		"device_type": deviceType,
		"clicked_at":  utils.NowUTC(),
	}

	if userID != "" {
		data["user_id"] = userID
	}

	resp, _, err := u.Client.
		From("analytics").
		Insert(data, false, "", "", "").
		Execute()

	if err != nil {
		log.Printf("[SaveAnalytics] Insert failed: %v", err)
		return fmt.Errorf("failed to save analytics: %w", err)
	}

	if len(resp) == 0 {
		log.Printf("[SaveAnalytics] Empty insert response for url_id=%s", urlID)
	} else {
		userIDPart := "anonymous"
		if userID != "" {
			userIDPart = userID
		}
		log.Printf("[SaveAnalytics] Successfully saved analytics for url_id=%s, user_id=%s", urlID, userIDPart)
	}

	return nil
}
