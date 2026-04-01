package repository

import (
	"context"
	"time"

	"url-shortener-go-backend/internal/metrics"
	"url-shortener-go-backend/internal/model"
)

type InstrumentedURLRepository struct {
	inner URLRepository
}

func NewInstrumentedURLRepository(inner URLRepository) URLRepository {
	return &InstrumentedURLRepository{inner: inner}
}

func (r *InstrumentedURLRepository) SaveURL(ctx context.Context, url *model.URL) error {
	start := time.Now()
	err := r.inner.SaveURL(ctx, url)
	metrics.DBQueryDuration.WithLabelValues("SaveURL", "urls").Observe(time.Since(start).Seconds())
	return err
}

func (r *InstrumentedURLRepository) GetURLByShortCode(ctx context.Context, shortcode string) (*model.URL, error) {
	start := time.Now()
	url, err := r.inner.GetURLByShortCode(ctx, shortcode)
	metrics.DBQueryDuration.WithLabelValues("GetURLByShortCode", "urls").Observe(time.Since(start).Seconds())
	return url, err
}

func (r *InstrumentedURLRepository) GetUserUrls(ctx context.Context, userID string) ([]model.URL, error) {
	start := time.Now()
	urls, err := r.inner.GetUserUrls(ctx, userID)
	metrics.DBQueryDuration.WithLabelValues("GetUserUrls", "urls").Observe(time.Since(start).Seconds())
	return urls, err
}

func (r *InstrumentedURLRepository) IncrementClickCount(ctx context.Context, shortcode string) error {
	start := time.Now()
	err := r.inner.IncrementClickCount(ctx, shortcode)
	metrics.DBQueryDuration.WithLabelValues("IncrementClickCount", "urls").Observe(time.Since(start).Seconds())
	return err
}

type InstrumentedAnalyticsRepository struct {
	inner AnalyticsRepository
}

func NewInstrumentedAnalyticsRepository(inner AnalyticsRepository) AnalyticsRepository {
	return &InstrumentedAnalyticsRepository{inner: inner}
}

func (r *InstrumentedAnalyticsRepository) SaveAnalytics(ctx context.Context, userID, urlID, referrer, deviceType string) error {
	start := time.Now()
	err := r.inner.SaveAnalytics(ctx, userID, urlID, referrer, deviceType)
	metrics.DBQueryDuration.WithLabelValues("SaveAnalytics", "analytics").Observe(time.Since(start).Seconds())
	return err
}

func (r *InstrumentedAnalyticsRepository) GetUserAnalyticsSummary(ctx context.Context, userID string) (*model.UserAnalyticsSummary, error) {
	start := time.Now()
	summary, err := r.inner.GetUserAnalyticsSummary(ctx, userID)
	metrics.DBQueryDuration.WithLabelValues("GetUserAnalyticsSummary", "analytics").Observe(time.Since(start).Seconds())
	return summary, err
}

func (r *InstrumentedAnalyticsRepository) GetUserTopURLs(ctx context.Context, userID string, limit int) ([]model.URLClickStats, error) {
	start := time.Now()
	urls, err := r.inner.GetUserTopURLs(ctx, userID, limit)
	metrics.DBQueryDuration.WithLabelValues("GetUserTopURLs", "urls").Observe(time.Since(start).Seconds())
	return urls, err
}

func (r *InstrumentedAnalyticsRepository) GetUserDailyClicks(ctx context.Context, userID string, days int) ([]model.DailyClickStats, error) {
	start := time.Now()
	stats, err := r.inner.GetUserDailyClicks(ctx, userID, days)
	metrics.DBQueryDuration.WithLabelValues("GetUserDailyClicks", "analytics").Observe(time.Since(start).Seconds())
	return stats, err
}

func (r *InstrumentedAnalyticsRepository) GetUserTopReferrers(ctx context.Context, userID string, limit int) ([]model.ReferrerStats, error) {
	start := time.Now()
	refs, err := r.inner.GetUserTopReferrers(ctx, userID, limit)
	metrics.DBQueryDuration.WithLabelValues("GetUserTopReferrers", "analytics").Observe(time.Since(start).Seconds())
	return refs, err
}

func (r *InstrumentedAnalyticsRepository) GetUserDeviceBreakdown(ctx context.Context, userID string) ([]model.DeviceStats, error) {
	start := time.Now()
	devices, err := r.inner.GetUserDeviceBreakdown(ctx, userID)
	metrics.DBQueryDuration.WithLabelValues("GetUserDeviceBreakdown", "analytics").Observe(time.Since(start).Seconds())
	return devices, err
}

func (r *InstrumentedAnalyticsRepository) AggregateYesterdayAnalytics(ctx context.Context) error {
	start := time.Now()
	err := r.inner.AggregateYesterdayAnalytics(ctx)
	metrics.DBQueryDuration.WithLabelValues("AggregateYesterdayAnalytics", "analytics").Observe(time.Since(start).Seconds())
	return err
}

func (r *InstrumentedAnalyticsRepository) GetUserStats(ctx context.Context, userID string) (int64, int64, int64, int64, error) {
	start := time.Now()
	a, b, c, d, err := r.inner.GetUserStats(ctx, userID)
	metrics.DBQueryDuration.WithLabelValues("GetUserStats", "analytics").Observe(time.Since(start).Seconds())
	return a, b, c, d, err
}
