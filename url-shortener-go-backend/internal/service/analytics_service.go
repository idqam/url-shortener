package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"url-shortener-go-backend/internal/cache"
	"url-shortener-go-backend/internal/metrics"
	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/repository"
)

type AnalyticsServiceImpl struct {
	analyticsRepo repository.AnalyticsRepository
	cache         cache.Cache
	salt          string
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository, c cache.Cache, salt string) AnalyticsService {
	return &AnalyticsServiceImpl{
		analyticsRepo: analyticsRepo,
		cache:         c,
		salt:          salt,
	}
}

func NormalizeSummary(summary *model.UserAnalyticsSummary) {
	if summary.TopURLs == nil {
		summary.TopURLs = []model.URLClickStats{}
	}
	if summary.TopReferrers == nil {
		summary.TopReferrers = []model.ReferrerStats{}
	}
	if summary.DeviceBreakdown == nil {
		summary.DeviceBreakdown = []model.DeviceStats{}
	}
	if summary.DailyClickTrend == nil {
		summary.DailyClickTrend = []model.DailyClickStats{}
	}
}

func (s *AnalyticsServiceImpl) GetUserDashboard(ctx context.Context, userID string) (*model.UserAnalyticsSummary, error) {
	cacheKey := cache.KeyUserAnalytics(s.salt, userID, time.Now().AddDate(0, 0, -7), time.Now())

	if val, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var cached model.UserAnalyticsSummary
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			slog.Info("dashboard cache hit", "user_id", userID)
			NormalizeSummary(&cached)
			return &cached, nil
		}
	}

	slog.Info("dashboard cache miss", "user_id", userID)
	summary, err := s.analyticsRepo.GetUserAnalyticsSummary(ctx, userID)
	if summary != nil {
		NormalizeSummary(summary)
	}

	if err != nil {
		slog.Error("failed to get analytics summary", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get user dashboard: %w", err)
	}

	if jsonVal, err := json.Marshal(summary); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(jsonVal), time.Hour)
		slog.Info("dashboard cached", "user_id", userID)
	}

	return summary, nil
}

func (s *AnalyticsServiceImpl) GetUserTopURLs(ctx context.Context, userID string, limit int) ([]model.URLClickStats, error) {
	cacheKey := fmt.Sprintf("user_top_urls:%s:%d", userID, limit)

	if val, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var urls []model.URLClickStats
		if err := json.Unmarshal([]byte(val), &urls); err == nil {
			return urls, nil
		}
	}

	urls, err := s.analyticsRepo.GetUserTopURLs(ctx, userID, limit)
	if err != nil {
		slog.Error("failed to get top urls", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get top URLs: %w", err)
	}

	if jsonVal, err := json.Marshal(urls); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(jsonVal), 30*time.Minute)
	}

	return urls, nil
}

func (s *AnalyticsServiceImpl) GetUserDailyTrend(ctx context.Context, userID string, days int) ([]model.DailyClickStats, error) {
	cacheKey := fmt.Sprintf("user_daily_trend:%s:%d", userID, days)

	if val, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var trend []model.DailyClickStats
		if err := json.Unmarshal([]byte(val), &trend); err == nil {
			return trend, nil
		}
	}

	trend, err := s.analyticsRepo.GetUserDailyClicks(ctx, userID, days)
	if err != nil {
		slog.Error("failed to get daily trend", "user_id", userID, "error", err)
		return []model.DailyClickStats{}, nil
	}

	if trend == nil {
		trend = []model.DailyClickStats{}
	}

	if jsonVal, err := json.Marshal(trend); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(jsonVal), 15*time.Minute)
	}

	return trend, nil
}

func (s *AnalyticsServiceImpl) GetUserTopReferrers(ctx context.Context, userID string, limit int) ([]model.ReferrerStats, error) {
	cacheKey := fmt.Sprintf("user_top_referrers:%s:%d", userID, limit)

	if val, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var referrers []model.ReferrerStats
		if err := json.Unmarshal([]byte(val), &referrers); err == nil {
			return referrers, nil
		}
	}

	referrers, err := s.analyticsRepo.GetUserTopReferrers(ctx, userID, limit)
	if err != nil {
		slog.Error("failed to get top referrers", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get top referrers: %w", err)
	}

	if jsonVal, err := json.Marshal(referrers); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(jsonVal), 45*time.Minute)
	}

	return referrers, nil
}

func (s *AnalyticsServiceImpl) GetUserDeviceBreakdown(ctx context.Context, userID string) ([]model.DeviceStats, error) {
	cacheKey := fmt.Sprintf("user_device_breakdown:%s", userID)

	if val, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var devices []model.DeviceStats
		if err := json.Unmarshal([]byte(val), &devices); err == nil {
			return devices, nil
		}
	}

	devices, err := s.analyticsRepo.GetUserDeviceBreakdown(ctx, userID)
	if err != nil {
		slog.Error("failed to get device breakdown", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get device breakdown: %w", err)
	}

	if jsonVal, err := json.Marshal(devices); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(jsonVal), time.Hour)
	}

	return devices, nil
}

func (s *AnalyticsServiceImpl) RecordAnalytics(ctx context.Context, userID, urlID, referrer, deviceType string) error {
	go func() {
		bgCtx := context.Background()
		if err := s.analyticsRepo.SaveAnalytics(bgCtx, userID, urlID, referrer, deviceType); err != nil {
			slog.Error("failed to save analytics", "user_id", userID, "url_id", urlID, "error", err)
		} else {
			metrics.AnalyticsRecordsTotal.Inc()
			s.invalidateUserCaches(bgCtx, userID)
		}
	}()

	return nil
}

func (s *AnalyticsServiceImpl) ProcessDailyAnalytics(ctx context.Context) error {
	slog.Info("starting daily analytics aggregation")

	if err := s.analyticsRepo.AggregateYesterdayAnalytics(ctx); err != nil {
		slog.Error("failed to aggregate analytics", "error", err)
		return fmt.Errorf("failed to process daily analytics: %w", err)
	}

	slog.Info("daily analytics aggregation completed")
	return nil
}

func (s *AnalyticsServiceImpl) invalidateUserCaches(ctx context.Context, userID string) {
	keysToDelete := []string{
		fmt.Sprintf("user_top_urls:%s:10", userID),
		fmt.Sprintf("user_daily_trend:%s:7", userID),
		fmt.Sprintf("user_top_referrers:%s:5", userID),
		fmt.Sprintf("user_device_breakdown:%s", userID),
		cache.KeyUserAnalytics(s.salt, userID, time.Now().AddDate(0, 0, -7), time.Now()),
	}

	for _, key := range keysToDelete {
		if err := s.cache.Delete(ctx, key); err != nil {
			slog.Warn("failed to delete cache key", "key", key, "error", err)
		}
	}
}
