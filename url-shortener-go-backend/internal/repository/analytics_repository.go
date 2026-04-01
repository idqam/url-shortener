package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"time"
	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/utils"

	"github.com/supabase-community/postgrest-go"
)

type AnalyticsRepositoryImpl struct {
	*SupabaseRepository
}

func NewAnalyticsRepository(baseRepo *SupabaseRepository) AnalyticsRepository {
	return &AnalyticsRepositoryImpl{baseRepo}
}

func (a *AnalyticsRepositoryImpl) SaveAnalytics(ctx context.Context, userID, urlID, referrer, deviceType string) error {
	data := map[string]interface{}{
		"url_id":      urlID,
		"referrer":    referrer,
		"device_type": deviceType,
		"clicked_at":  utils.NowUTC(),
	}

	if userID != "" {
		data["user_id"] = userID
	}

	resp, _, err := a.Client.
		From("analytics").
		Insert(data, false, "", "", "").
		Execute()

	if err != nil {
		slog.Error("analytics insert failed", "error", err)
		return fmt.Errorf("failed to save analytics: %w", err)
	}

	if len(resp) == 0 {
		slog.Warn("analytics insert returned empty response", "url_id", urlID)
	} else {
		userIDPart := "anonymous"
		if userID != "" {
			userIDPart = userID
		}
		slog.Info("analytics saved", "url_id", urlID, "user_id", userIDPart)
	}

	return nil
}

func (a *AnalyticsRepositoryImpl) GetUserAnalyticsSummary(ctx context.Context, userID string) (*model.UserAnalyticsSummary, error) {
	slog.Info("creating analytics summary", "user_id", userID)

	summary := &model.UserAnalyticsSummary{
		TopURLs:         []model.URLClickStats{},
		TopReferrers:    []model.ReferrerStats{},
		DeviceBreakdown: []model.DeviceStats{},
		DailyClickTrend: []model.DailyClickStats{},
	}

	topURLs, err := a.GetUserTopURLs(ctx, userID, 100)
	if err != nil {
		slog.Error("failed to get top urls for summary", "user_id", userID, "error", err)
		summary.TopURLs = []model.URLClickStats{}
		summary.TotalURLs = 0
		summary.TotalClicks = 0
		summary.AverageClicks = 0.0
	} else {
		summary.TotalURLs = int64(len(topURLs))
		summary.TotalClicks = 0
		for _, url := range topURLs {
			summary.TotalClicks += url.ClickCount
		}

		if summary.TotalURLs > 0 {
			summary.AverageClicks = float64(summary.TotalClicks) / float64(summary.TotalURLs)
		} else {
			summary.AverageClicks = 0.0
		}

		if len(topURLs) > 10 {
			summary.TopURLs = topURLs[:10]
		} else {
			summary.TopURLs = topURLs
		}

		slog.Info("calculated url stats", "user_id", userID, "total_urls", summary.TotalURLs, "total_clicks", summary.TotalClicks)
	}

	dailyTrend, err := a.GetUserDailyClicks(ctx, userID, 2)
	if err != nil || len(dailyTrend) == 0 {
		slog.Warn("could not get daily trend for today/yesterday", "user_id", userID)
		summary.ClicksToday = 0
		summary.ClicksYesterday = 0
	} else {
		today := time.Now().Format("2006-01-02")
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

		summary.ClicksToday = 0
		summary.ClicksYesterday = 0

		for _, day := range dailyTrend {
			switch day.Date {
			case today:
				summary.ClicksToday = day.Clicks
			case yesterday:
				summary.ClicksYesterday = day.Clicks
			}
		}

		slog.Info("daily clicks calculated", "user_id", userID, "today", summary.ClicksToday, "yesterday", summary.ClicksYesterday)
	}

	topReferrers, err := a.GetUserTopReferrers(ctx, userID, 5)
	if err != nil {
		slog.Error("failed to get top referrers for summary", "user_id", userID, "error", err)
		summary.TopReferrers = []model.ReferrerStats{}
	} else {
		summary.TopReferrers = topReferrers
	}

	deviceBreakdown, err := a.GetUserDeviceBreakdown(ctx, userID)
	if err != nil {
		slog.Error("failed to get device breakdown for summary", "user_id", userID, "error", err)
		summary.DeviceBreakdown = []model.DeviceStats{}
	} else {
		summary.DeviceBreakdown = deviceBreakdown
	}

	dailyTrend, err = a.GetUserDailyClicks(ctx, userID, 7)
	if err != nil {
		slog.Error("failed to get daily trend for summary", "user_id", userID, "error", err)
		summary.DailyClickTrend = []model.DailyClickStats{}
	} else {
		summary.DailyClickTrend = dailyTrend
	}

	slog.Info("analytics summary created", "user_id", userID)
	return summary, nil
}

func (a *AnalyticsRepositoryImpl) GetUserStats(ctx context.Context, userID string) (totalURLs int64, totalClicks int64, clicksToday int64, clicksYesterday int64, err error) {
	resp, count, err := a.Client.
		From("urls").
		Select("click_count", "exact", true).
		Eq("user_id", userID).
		Execute()

	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to get user URL stats: %w", err)
	}

	totalURLs = int64(count)

	var urls []struct {
		ClickCount int64 `json:"click_count"`
	}
	if err := json.Unmarshal(resp, &urls); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to decode URL stats: %w", err)
	}

	for _, url := range urls {
		totalClicks += url.ClickCount
	}

	today := time.Now().Format("2006-01-02")
	_, todayCount, err := a.Client.
		From("analytics").
		Select("id", "exact", true).
		Eq("user_id", userID).
		Gte("clicked_at", today).
		Execute()

	if err != nil {
		slog.Error("failed to get today's clicks", "user_id", userID, "error", err)
		clicksToday = 0
	} else {
		clicksToday = int64(todayCount)
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	yesterdayEnd := today
	_, yesterdayCount, err := a.Client.
		From("analytics").
		Select("id", "exact", true).
		Eq("user_id", userID).
		Gte("clicked_at", yesterday).
		Lt("clicked_at", yesterdayEnd).
		Execute()

	if err != nil {
		slog.Error("failed to get yesterday's clicks", "user_id", userID, "error", err)
		clicksYesterday = 0
	} else {
		clicksYesterday = int64(yesterdayCount)
	}

	return totalURLs, totalClicks, clicksToday, clicksYesterday, nil
}

func (a *AnalyticsRepositoryImpl) GetUserTopURLs(ctx context.Context, userID string, limit int) ([]model.URLClickStats, error) {
	resp, _, err := a.Client.
		From("urls").
		Select("id, short_code, original_url, click_count, created_at", "exact", false).
		Eq("user_id", userID).
		Order("click_count", &postgrest.OrderOpts{Ascending: false}).
		Limit(limit, "").
		Execute()

	if err != nil {
		return []model.URLClickStats{}, fmt.Errorf("failed to fetch top URLs: %w", err)
	}

	var urls []model.URLClickStats
	if err := json.Unmarshal(resp, &urls); err != nil {
		return []model.URLClickStats{}, fmt.Errorf("failed to decode top URLs: %w", err)
	}

	if urls == nil {
		urls = []model.URLClickStats{}
	}

	return urls, nil
}

func fillMissingDays(data []model.DailyClickStats, days int) []model.DailyClickStats {
	today := time.Now().UTC()
	clicksByDate := make(map[string]int64)
	for _, d := range data {
		clicksByDate[d.Date] = d.Clicks
	}

	result := make([]model.DailyClickStats, days)
	for i := 0; i < days; i++ {
		date := today.AddDate(0, 0, -(days-1-i)).Format("2006-01-02")
		result[i] = model.DailyClickStats{
			Date:   date,
			Clicks: clicksByDate[date],
		}
	}
	return result
}

func (a *AnalyticsRepositoryImpl) GetUserDailyClicks(ctx context.Context, userID string, days int) ([]model.DailyClickStats, error) {
	rawJSON := a.Client.Rpc("get_user_daily_clicks", "", map[string]any{
		"p_days":    days,
		"p_user_id": userID,
	})

	if rawJSON == "" {
		return []model.DailyClickStats{}, nil
	}

	var stats []model.DailyClickStats
	if err := json.Unmarshal([]byte(rawJSON), &stats); err != nil {
		return []model.DailyClickStats{}, fmt.Errorf("failed to unmarshal daily clicks: %w", err)
	}

	stats = fillMissingDays(stats, days)

	return stats, nil
}

func (a *AnalyticsRepositoryImpl) GetUserTopReferrers(ctx context.Context, userID string, limit int) ([]model.ReferrerStats, error) {
	resp, _, err := a.Client.
		From("analytics").
		Select("referrer", "exact", false).
		Eq("user_id", userID).
		Not("referrer", "is", "").
		Not("referrer", "eq", "").
		Execute()

	if err != nil {
		return []model.ReferrerStats{}, fmt.Errorf("failed to fetch referrer data: %w", err)
	}

	if len(resp) == 0 || string(resp) == "[]" || string(resp) == "" {
		slog.Info("no referrer data found", "user_id", userID)
		return []model.ReferrerStats{}, nil
	}

	var analytics []struct {
		Referrer string `json:"referrer"`
	}
	if err := json.Unmarshal(resp, &analytics); err != nil {
		slog.Error("failed to decode referrer data", "user_id", userID, "error", err)
		return []model.ReferrerStats{}, nil
	}

	referrerCounts := make(map[string]int64)
	for _, record := range analytics {
		if record.Referrer != "" {
			referrerCounts[record.Referrer]++
		}
	}

	referrers := make([]model.ReferrerStats, 0, len(referrerCounts))
	for referrer, count := range referrerCounts {
		referrers = append(referrers, model.ReferrerStats{
			Referrer: referrer,
			Clicks:   count,
		})
	}

	slices.SortFunc(referrers, func(a, b model.ReferrerStats) int {
		if a.Clicks > b.Clicks {
			return -1
		}
		if a.Clicks < b.Clicks {
			return 1
		}
		return 0
	})

	if len(referrers) > limit {
		referrers = referrers[:limit]
	}

	slog.Info("referrers found", "user_id", userID, "count", len(referrers))
	return referrers, nil
}

func (a *AnalyticsRepositoryImpl) GetUserDeviceBreakdown(ctx context.Context, userID string) ([]model.DeviceStats, error) {
	resp, _, err := a.Client.
		From("daily_analytics").
		Select(`
			'desktop' as device_type, SUM(desktop_clicks) as clicks
			UNION ALL
			SELECT 'mobile' as device_type, SUM(mobile_clicks) as clicks
			UNION ALL
			SELECT 'tablet' as device_type, SUM(tablet_clicks) as clicks
			UNION ALL
			SELECT 'unknown' as device_type, SUM(unknown_clicks) as clicks
		`, "exact", false).
		Eq("user_id", userID).
		Order("clicks", &postgrest.OrderOpts{Ascending: false}).
		Execute()

	if err != nil {
		return a.getUserDeviceBreakdownFromRaw(ctx, userID)
	}

	var devices []model.DeviceStats
	if err := json.Unmarshal(resp, &devices); err != nil {
		return a.getUserDeviceBreakdownFromRaw(ctx, userID)
	}

	if devices == nil {
		devices = []model.DeviceStats{}
	}

	return devices, nil
}

func (a *AnalyticsRepositoryImpl) getUserDeviceBreakdownFromRaw(ctx context.Context, userID string) ([]model.DeviceStats, error) {
	resp, _, err := a.Client.
		From("analytics").
		Select("device_type", "exact", false).
		Eq("user_id", userID).
		Execute()

	if err != nil {
		return []model.DeviceStats{}, fmt.Errorf("failed to fetch device data: %w", err)
	}

	var analytics []struct {
		DeviceType *string `json:"device_type"`
	}
	if err := json.Unmarshal(resp, &analytics); err != nil {
		return []model.DeviceStats{}, fmt.Errorf("failed to decode device data: %w", err)
	}

	deviceCounts := make(map[string]int64)
	for _, record := range analytics {
		deviceType := "unknown"
		if record.DeviceType != nil && *record.DeviceType != "" {
			deviceType = *record.DeviceType
		}
		deviceCounts[deviceType]++
	}

	devices := make([]model.DeviceStats, 0, len(deviceCounts))
	for deviceType, count := range deviceCounts {
		devices = append(devices, model.DeviceStats{
			DeviceType: deviceType,
			Clicks:     count,
		})
	}

	slices.SortFunc(devices, func(a, b model.DeviceStats) int {
		if a.Clicks > b.Clicks {
			return -1
		}
		if a.Clicks < b.Clicks {
			return 1
		}
		return 0
	})

	return devices, nil
}

func (a *AnalyticsRepositoryImpl) AggregateYesterdayAnalytics(ctx context.Context) error {
	err := a.Client.Rpc("update_daily_analytics", "", map[string]any{})
	if err != "" {
		slog.Error("rpc update_daily_analytics failed", "error", err)
		return fmt.Errorf("failed to aggregate analytics: %s", err)
	}

	slog.Info("yesterday analytics aggregated successfully")
	return nil
}
