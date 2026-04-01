package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/utils"
)

type URLRepositoryImpl struct {
	*SupabaseRepository
	shortDomain string
}

func NewURLRepository(baseRepo *SupabaseRepository, shortDomain string) URLRepository {
	return &URLRepositoryImpl{SupabaseRepository: baseRepo, shortDomain: shortDomain}
}

func (u *URLRepositoryImpl) GetURLByShortCode(ctx context.Context, shortcode string) (*model.URL, error) {
	resp, _, err := u.Client.
		From("urls").
		Select("id, original_url, short_code, click_count, is_public, created_at", "exact", false).
		Eq("short_code", shortcode).
		Single().
		Execute()

	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "pgrst116") || strings.Contains(errStr, "no rows") || strings.Contains(errStr, "json object requested") {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}

	if len(resp) == 0 || string(resp) == "null" {
		return nil, utils.ErrNotFound
	}

	var url model.URL
	if err := json.Unmarshal(resp, &url); err != nil {
		return nil, fmt.Errorf("failed to decode URL response: %w", err)
	}

	url.PopulateShortURL(u.shortDomain)

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

	if urls == nil {
		urls = []model.URL{}
	}

	for i := range urls {
		urls[i].PopulateShortURL(u.shortDomain)
	}

	return urls, nil
}

func (u *URLRepositoryImpl) IncrementClickCount(ctx context.Context, shortcode string) error {
	err := u.Client.Rpc("increment_click_count", "", map[string]any{
		"sc": shortcode,
	})

	if err != "" {
		slog.Error("rpc increment_click_count failed", "shortcode", shortcode, "error", err)
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
		slog.Error("url insert failed", "error", err)
		return fmt.Errorf("supabase insert failed: %w", err)
	}

	if len(resp) == 0 {
		slog.Error("url insert returned empty response")
		return fmt.Errorf("failed to save URL (empty response)")
	}

	var inserted []model.URL
	if err := json.Unmarshal(resp, &inserted); err != nil {
		var errResp struct {
			Message string `json:"message"`
			Code    string `json:"code"`
			Details string `json:"details"`
			Hint    string `json:"hint"`
		}
		if json.Unmarshal(resp, &errResp) == nil && errResp.Message != "" {
			msg := strings.ToLower(errResp.Message)
			if strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint") {
				slog.Warn("unique constraint violation on url insert", "message", errResp.Message)
				return ErrUniqueViolation
			}
			return fmt.Errorf("supabase error: %s", errResp.Message)
		}

		return fmt.Errorf("failed to decode inserted URL: %s", string(resp))
	}

	if len(inserted) == 0 {
		slog.Error("url insert returned no rows")
		return fmt.Errorf("no URL returned after insert")
	}

	*url = inserted[0]

	url.PopulateShortURL(u.shortDomain)
	slog.Info("url insert successful", "id", url.ID, "short_code", url.ShortCode)
	return nil
}
