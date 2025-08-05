package repository

import (
	"context"
	"url-shortener-go-backend/internal/model"
)

type URLRepository interface {
	SaveURL(ctx context.Context, url *model.URL) error
	GetURLByShortCode(ctx context.Context, shortcode string) (*model.URL, error)
	GetUserUrls(ctx context.Context, userID string) ([]model.URL, error)
	IncrementClickCount(ctx context.Context, shortcode string) error
	SaveAnalytics(ctx context.Context, userID, urlID, referrer, deviceType string) error
}
