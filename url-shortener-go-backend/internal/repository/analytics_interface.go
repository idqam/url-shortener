package repository

import (
	"context"
	"url-shortener-go-backend/internal/model"
)

type AnalyticsRepository interface {
	SaveAnalytics(ctx context.Context, userID, urlID, referrer, deviceType string) error

	GetUserAnalyticsSummary(ctx context.Context, userID string) (*model.UserAnalyticsSummary, error)

	GetUserTopURLs(ctx context.Context, userID string, limit int) ([]model.URLClickStats, error)

	GetUserDailyClicks(ctx context.Context, userID string, days int) ([]model.DailyClickStats, error)

	GetUserTopReferrers(ctx context.Context, userID string, limit int) ([]model.ReferrerStats, error)

	GetUserDeviceBreakdown(ctx context.Context, userID string) ([]model.DeviceStats, error)

	AggregateYesterdayAnalytics(ctx context.Context) error

	GetUserStats(ctx context.Context, userID string) (totalURLs int64, totalClicks int64, clicksToday int64, clicksYesterday int64, err error)
}
