package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"url-shortener-go-backend/internal/cache"
	"url-shortener-go-backend/internal/metrics"
	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/repository"
	"url-shortener-go-backend/internal/utils"
)

type URLService interface {
	CreateShortURL(ctx context.Context, originalURL string, isPublic bool, userID *string, codeLength int) (*model.URL, error)
	GetURLByShortCode(ctx context.Context, shortcode string) (*model.URL, error)
	GetUserUrls(ctx context.Context, userID string) ([]model.URL, error)
	IncrementClickCount(ctx context.Context, shortcode string) error
}

type URLServiceImpl struct {
	repo  repository.URLRepository
	cache cache.Cache
	salt  string
}

func NewURLService(repo repository.URLRepository, c cache.Cache, salt string) URLService {
	return &URLServiceImpl{
		repo:  repo,
		cache: c,
		salt:  salt,
	}
}

func (s *URLServiceImpl) CreateShortURL(ctx context.Context, originalURL string, isPublic bool, userID *string, codeLength int) (*model.URL, error) {
	cacheKey := fmt.Sprintf("short_url:%s:%v", originalURL, userID)
	if val, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var cachedURL model.URL
		if err := json.Unmarshal([]byte(val), &cachedURL); err == nil {
			return &cachedURL, nil
		}
	}

	shortcode, err := utils.GenerateCode(originalURL, codeLength, s.salt)
	if err != nil {
		slog.Error("failed to generate shortcode", "error", err)
		return nil, fmt.Errorf("failed to generate shortcode")
	}

	url := &model.URL{
		ShortCode:   shortcode,
		OriginalURL: originalURL,
		IsPublic:    isPublic,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	for retries := 0; retries < 3; retries++ {
		err := s.repo.SaveURL(ctx, url)
		if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			shortcode, err := utils.GenerateCode(originalURL, codeLength, s.salt)
			if err != nil {
				slog.Error("failed to generate shortcode on retry", "error", err)
				return nil, fmt.Errorf("failed to generate shortcode")
			}
			url.ShortCode = shortcode
			continue
		}
		if err != nil {
			slog.Error("failed to save url", "error", err)
			return nil, err
		}
		break
	}

	metrics.URLShortensTotal.Inc()

	if jsonVal, err := json.Marshal(url); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(jsonVal), time.Hour)
	}

	return url, nil
}

func (s *URLServiceImpl) GetURLByShortCode(ctx context.Context, shortcode string) (*model.URL, error) {
	cacheKey := "short_url:" + shortcode

	if val, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var url model.URL
		if err := json.Unmarshal([]byte(val), &url); err == nil {
			return &url, nil
		}
	}

	url, err := s.repo.GetURLByShortCode(ctx, shortcode)
	if err != nil {
		slog.Error("db lookup failed for shortcode", "shortcode", shortcode, "error", err)
		return nil, err
	}

	jsonVal, _ := json.Marshal(url)
	_ = s.cache.Set(ctx, cacheKey, string(jsonVal), time.Hour)

	return url, nil
}

func (s *URLServiceImpl) GetUserUrls(ctx context.Context, userID string) ([]model.URL, error) {
	cacheKey := "user_urls:" + userID

	if val, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var urls []model.URL
		if err := json.Unmarshal([]byte(val), &urls); err == nil {
			return urls, nil
		}
	}

	urls, err := s.repo.GetUserUrls(ctx, userID)
	if err != nil {
		slog.Error("db fetch failed for user urls", "user_id", userID, "error", err)
		return nil, err
	}

	jsonVal, _ := json.Marshal(urls)
	_ = s.cache.Set(ctx, cacheKey, string(jsonVal), time.Hour)

	return urls, nil
}

func (s *URLServiceImpl) IncrementClickCount(ctx context.Context, shortcode string) error {
	err := s.repo.IncrementClickCount(ctx, shortcode)
	if err != nil {
		slog.Error("failed to increment click count", "shortcode", shortcode, "error", err)
	}
	return err
}
