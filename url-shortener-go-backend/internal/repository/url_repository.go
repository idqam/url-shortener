package repository

import (
	"encoding/json"
	"fmt"
	"strings"
	"url-shortener-go-backend/internal/model"
)

type URLRepositoryImpl struct {
	*SupabaseRepository
}
func (u *URLRepositoryImpl) GetURLByShortCode(shortcode string) (*model.URL, error) {
    resp, _, err := u.Client.
        From("urls").
        Select("*", "exact", false).
        Eq("short_code", shortcode).
        Execute()


    if err != nil {
        return nil, fmt.Errorf("failed to fetch URL by shortcode: %w", err)
    }
    if len(resp) == 0 || string(resp) == "[]" {
        return nil, nil
    }

    var results []model.URL
    if err := json.Unmarshal(resp, &results); err != nil {
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
		if strings.Contains(err.Error(), "No rows") || strings.Contains(strings.ToLower(err.Error()), "not found") {
			return []model.URL{}, nil
		}
		return nil, fmt.Errorf("failed to get user urls: %w", err)
	}

	var urls []model.URL
	if err := json.Unmarshal(resp, &urls); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return urls, nil
}

func (u *URLRepositoryImpl) IncrementClickCount(shortcode string) error {
	err := u.Client.Rpc("increment_click_amount", "", map[string]any{
		"sc": shortcode,
	})

	if err != "" {
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
		"created_at":   url.CreatedAt,
	}

	if url.UserID != nil && *url.UserID != "" {
		data["user_id"] = *url.UserID
	}

	_, _, err := u.Client.
		From("urls").
		Insert(data, false, "", "return=minimal", "").
		Execute()


	if err != nil {
		return fmt.Errorf("failed to save URL: %w", err)
	}
	return nil
}





func NewURLRepository(baseRepo *SupabaseRepository) URLRepository {
	return &URLRepositoryImpl{baseRepo}
}
