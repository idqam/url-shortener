package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"url-shortener-go-backend/internal/cache"
	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/repository"
)

type AnalyticsServiceImpl struct {
	analyticsRepo repository.AnalyticsRepository
	cache         cache.Cache
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository, cache cache.Cache) AnalyticsService {
	return &AnalyticsServiceImpl{
		analyticsRepo: analyticsRepo,
		cache:         cache,
	}
}

func isValidSummary(summary *model.UserAnalyticsSummary) bool {
	return summary != nil &&
		summary.DailyClickTrend != nil &&
		summary.DeviceBreakdown != nil &&
		summary.TopReferrers != nil &&
		summary.TopURLs != nil
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
    cacheKey := cache.KeyUserAnalytics(userID, time.Now().AddDate(0, 0, -7), time.Now())

    if val, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
        var summary model.UserAnalyticsSummary
        if err := json.Unmarshal([]byte(val), &summary); err == nil {
            log.Printf("[GetUserDashboard] Cache HIT for user: %s", userID)
            NormalizeSummary(&summary) 
			 if jsonVal, err := json.Marshal(&summary); err == nil {
            _ = s.cache.Set(ctx, cacheKey, string(jsonVal), time.Hour)
        }
            return &summary, nil
        }
    }

    log.Printf("[GetUserDashboard] Cache MISS for user: %s", userID)
    summary, err := s.analyticsRepo.GetUserAnalyticsSummary(ctx, userID)
    if summary != nil {
        NormalizeSummary(summary)
    }

    if err != nil {
        log.Printf("[GetUserDashboard] Failed to get analytics summary for user %s: %v", userID, err)
        return nil, fmt.Errorf("failed to get user dashboard: %w", err)
    }

    if jsonVal, err := json.Marshal(summary); err == nil {
        _ = s.cache.Set(ctx, cacheKey, string(jsonVal), time.Hour)
        log.Printf("[GetUserDashboard] Cached analytics for user: %s", userID)
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
		log.Printf("[GetUserTopURLs] Failed for user %s: %v", userID, err)
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
		log.Printf("[GetUserDailyTrend] Failed for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to get daily trend: %w", err)
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
		log.Printf("[GetUserTopReferrers] Failed for user %s: %v", userID, err)
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
		log.Printf("[GetUserDeviceBreakdown] Failed for user %s: %v", userID, err)
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
			log.Printf("[RecordAnalytics] Failed to save analytics: userID=%s, urlID=%s, error=%v", userID, urlID, err)
		} else {

			s.invalidateUserCaches(bgCtx, userID)
		}
	}()

	return nil
}

// Implement cron scheduler?
func (s *AnalyticsServiceImpl) ProcessDailyAnalytics(ctx context.Context) error {
	log.Println("[ProcessDailyAnalytics] Starting daily analytics aggregation")

	if err := s.analyticsRepo.AggregateYesterdayAnalytics(ctx); err != nil {
		log.Printf("[ProcessDailyAnalytics] Failed to aggregate analytics: %v", err)
		return fmt.Errorf("failed to process daily analytics: %w", err)
	}

	log.Println("[ProcessDailyAnalytics] Successfully completed daily analytics aggregation")
	return nil
}

// TODO: If this becomes a real service then properly implement this, Redis pattern-based deletion and proper cache invalidation
func (s *AnalyticsServiceImpl) invalidateUserCaches(ctx context.Context, userID string) {

	patterns := []string{
		fmt.Sprintf("user_top_urls:%s:*", userID),
		fmt.Sprintf("user_daily_trend:%s:*", userID),
		fmt.Sprintf("user_top_referrers:%s:*", userID),
		fmt.Sprintf("user_device_breakdown:%s", userID),
	}

	for _, pattern := range patterns {
		log.Printf("[InvalidateCache] Pattern: %s", pattern)
	}

	cacheKey := cache.KeyUserAnalytics(userID, time.Now().AddDate(0, 0, -7), time.Now())

	log.Printf("[InvalidateCache] Dashboard cache key: %s", cacheKey)
}
