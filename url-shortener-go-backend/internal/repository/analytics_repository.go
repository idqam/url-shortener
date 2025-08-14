package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		log.Printf("[SaveAnalytics] Insert failed: %v", err)
		return fmt.Errorf("failed to save analytics: %w", err)
	}

	if len(resp) == 0 {
		log.Printf("[SaveAnalytics] Empty insert response for url_id=%s", urlID)
	} else {
		userIDPart := "anonymous"
		if userID != "" {
			userIDPart = userID
		}
		log.Printf("[SaveAnalytics] Successfully saved analytics for url_id=%s, user_id=%s", urlID, userIDPart)
	}

	return nil
}

func (a *AnalyticsRepositoryImpl) GetUserAnalyticsSummary(ctx context.Context, userID string) (*model.UserAnalyticsSummary, error) {
	log.Printf("[GetUserAnalyticsSummary] Creating summary with CALCULATED stats for user: %s", userID)

	summary := &model.UserAnalyticsSummary{
		TopURLs:         []model.URLClickStats{},
		TopReferrers:    []model.ReferrerStats{},
		DeviceBreakdown: []model.DeviceStats{},
		DailyClickTrend: []model.DailyClickStats{},
	}

	topURLs, err := a.GetUserTopURLs(ctx, userID, 100)
	if err != nil {
		log.Printf("[GetUserAnalyticsSummary] Failed to get top URLs: %v", err)
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

		log.Printf("[GetUserAnalyticsSummary] CALCULATED stats: URLs=%d, TotalClicks=%d, Avg=%.1f",
			summary.TotalURLs, summary.TotalClicks, summary.AverageClicks)
	}

	dailyTrend, err := a.GetUserDailyClicks(ctx, userID, 2)
	if err != nil || len(dailyTrend) == 0 {
		log.Printf("[GetUserAnalyticsSummary] Could not get daily trend for today/yesterday calculation")
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

		log.Printf("[GetUserAnalyticsSummary] Daily clicks: Today=%d, Yesterday=%d",
			summary.ClicksToday, summary.ClicksYesterday)
	}

	topReferrers, err := a.GetUserTopReferrers(ctx, userID, 5)
	if err != nil {
		log.Printf("[GetUserAnalyticsSummary] Failed to get top referrers: %v", err)
		summary.TopReferrers = []model.ReferrerStats{}
	} else {
		summary.TopReferrers = topReferrers
		log.Printf("[GetUserAnalyticsSummary] Got %d top referrers", len(topReferrers))
	}

	deviceBreakdown, err := a.GetUserDeviceBreakdown(ctx, userID)
	if err != nil {
		log.Printf("[GetUserAnalyticsSummary] Failed to get device breakdown: %v", err)
		summary.DeviceBreakdown = []model.DeviceStats{}
	} else {
		summary.DeviceBreakdown = deviceBreakdown
		log.Printf("[GetUserAnalyticsSummary] Got %d device types", len(deviceBreakdown))
	}

	dailyTrend, err = a.GetUserDailyClicks(ctx, userID, 7)
	if err != nil {
		log.Printf("[GetUserAnalyticsSummary] Failed to get daily trend: %v", err)
		summary.DailyClickTrend = []model.DailyClickStats{}
	} else {
		summary.DailyClickTrend = dailyTrend
		log.Printf("[GetUserAnalyticsSummary] Got %d days of trend data", len(dailyTrend))
	}

	log.Printf("[GetUserAnalyticsSummary] ✅ Successfully created CALCULATED summary for user %s", userID)
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
	resp, todayCount, err := a.Client.
		From("analytics").
		Select("id", "exact", true).
		Eq("user_id", userID).
		Gte("clicked_at", today).
		Execute()

	if err != nil {
		log.Printf("[GetUserStats] Failed to get today's clicks: %v", err)
		clicksToday = 0
	} else {
		clicksToday = int64(todayCount)
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	yesterdayEnd := today
	resp, yesterdayCount, err := a.Client.
		From("analytics").
		Select("id", "exact", true).
		Eq("user_id", userID).
		Gte("clicked_at", yesterday).
		Lt("clicked_at", yesterdayEnd).
		Execute()

	if err != nil {
		log.Printf("[GetUserStats] Failed to get yesterday's clicks: %v", err)
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
		return nil, fmt.Errorf("failed to fetch top URLs: %w", err)
	}

	var urls []model.URLClickStats
	if err := json.Unmarshal(resp, &urls); err != nil {
		return nil, fmt.Errorf("failed to decode top URLs: %w", err)
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
		date := today.AddDate(0, 0, -(days - 1 - i)).Format("2006-01-02")
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
		return nil, fmt.Errorf("failed to unmarshal daily clicks: %w", err)
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
		return nil, fmt.Errorf("failed to fetch referrer data: %w", err)
	}

	if len(resp) == 0 || string(resp) == "[]" || string(resp) == "" {
		log.Printf("[GetUserTopReferrers] No referrer data found for user %s", userID)
		return []model.ReferrerStats{}, nil
	}

	var analytics []struct {
		Referrer string `json:"referrer"`
	}
	if err := json.Unmarshal(resp, &analytics); err != nil {
		log.Printf("[GetUserTopReferrers] Failed to decode referrer data: %v", err)
		return []model.ReferrerStats{}, nil
	}

	referrerCounts := make(map[string]int64)
	for _, record := range analytics {
		if record.Referrer != "" {
			referrerCounts[record.Referrer]++
		}
	}

	var referrers []model.ReferrerStats
	for referrer, count := range referrerCounts {
		referrers = append(referrers, model.ReferrerStats{
			Referrer: referrer,
			Clicks:   count,
		})
	}

	if len(referrers) > 1 {
		for i := 0; i < len(referrers)-1; i++ {
			for j := i + 1; j < len(referrers); j++ {
				if referrers[i].Clicks < referrers[j].Clicks {
					referrers[i], referrers[j] = referrers[j], referrers[i]
				}
			}
		}
	}

	if len(referrers) > limit {
		referrers = referrers[:limit]
	}

	log.Printf("[GetUserTopReferrers] Found %d referrers for user %s", len(referrers), userID)
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

	return devices, nil
}

func (a *AnalyticsRepositoryImpl) getUserDeviceBreakdownFromRaw(ctx context.Context, userID string) ([]model.DeviceStats, error) {
	resp, _, err := a.Client.
		From("analytics").
		Select("device_type", "exact", false).
		Eq("user_id", userID).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch device data: %w", err)
	}

	var analytics []struct {
		DeviceType *string `json:"device_type"`
	}
	if err := json.Unmarshal(resp, &analytics); err != nil {
		return nil, fmt.Errorf("failed to decode device data: %w", err)
	}

	deviceCounts := make(map[string]int64)
	for _, record := range analytics {
		deviceType := "unknown"
		if record.DeviceType != nil && *record.DeviceType != "" {
			deviceType = *record.DeviceType
		}
		deviceCounts[deviceType]++
	}

	var devices []model.DeviceStats
	for deviceType, count := range deviceCounts {
		devices = append(devices, model.DeviceStats{
			DeviceType: deviceType,
			Clicks:     count,
		})
	}

	if len(devices) > 1 {
		for i := 0; i < len(devices)-1; i++ {
			for j := i + 1; j < len(devices); j++ {
				if devices[i].Clicks < devices[j].Clicks {
					devices[i], devices[j] = devices[j], devices[i]
				}
			}
		}
	}

	return devices, nil
}

func (a *AnalyticsRepositoryImpl) AggregateYesterdayAnalytics(ctx context.Context) error {
	err := a.Client.Rpc("update_daily_analytics", "", map[string]any{})
	if err != "" {
		log.Printf("[AggregateYesterdayAnalytics] RPC failed: %s", err)
		return fmt.Errorf("failed to aggregate analytics: %s", err)
	}

	log.Printf("[AggregateYesterdayAnalytics] Successfully aggregated yesterday's analytics")
	return nil
}
