package service

import (
	"context"
	"url-shortener-go-backend/internal/model"
)

type AnalyticsService interface {
	GetUserDashboard(ctx context.Context, userID string) (*model.UserAnalyticsSummary, error)

	GetUserTopURLs(ctx context.Context, userID string, limit int) ([]model.URLClickStats, error)
	GetUserDailyTrend(ctx context.Context, userID string, days int) ([]model.DailyClickStats, error)
	GetUserTopReferrers(ctx context.Context, userID string, limit int) ([]model.ReferrerStats, error)
	GetUserDeviceBreakdown(ctx context.Context, userID string) ([]model.DeviceStats, error)

	RecordAnalytics(ctx context.Context, userID, urlID, referrer, deviceType string) error

	ProcessDailyAnalytics(ctx context.Context) error
}
