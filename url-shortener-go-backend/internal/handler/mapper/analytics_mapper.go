package mapper

import (
	"url-shortener-go-backend/internal/handler/dto"
	"url-shortener-go-backend/internal/model"
)


func ToAnalyticsDashboardResponse(summary model.UserAnalyticsSummary) dto.AnalyticsDashboardResponse {
    return dto.AnalyticsDashboardResponse{
        Overview:        ToAnalyticsOverview(summary),
        TopURLs:         ToTopURLResponses(orEmptyURLs(summary.TopURLs)),
        TopReferrers:    ToReferrerResponses(orEmptyReferrers(summary.TopReferrers)),
        DeviceBreakdown: ToDeviceResponses(orEmptyDevices(summary.DeviceBreakdown)),
        DailyTrend:      ToDailyTrendResponses(orEmptyDailyTrend(summary.DailyClickTrend)),
    }
}

func orEmptyURLs(v []model.URLClickStats) []model.URLClickStats {
    if v == nil { return []model.URLClickStats{} }
    return v
}
func orEmptyReferrers(v []model.ReferrerStats) []model.ReferrerStats {
    if v == nil { return []model.ReferrerStats{} }
    return v
}
func orEmptyDevices(v []model.DeviceStats) []model.DeviceStats {
    if v == nil { return []model.DeviceStats{} }
    return v
}
func orEmptyDailyTrend(v []model.DailyClickStats) []model.DailyClickStats {
    if v == nil { return []model.DailyClickStats{} }
    return v
}


func ToAnalyticsOverview(summary model.UserAnalyticsSummary) dto.AnalyticsOverview {
	trendDirection := "same"
	if summary.ClicksToday > summary.ClicksYesterday {
		trendDirection = "up"
	} else if summary.ClicksToday < summary.ClicksYesterday {
		trendDirection = "down"
	}

	return dto.AnalyticsOverview{
		TotalURLs:       summary.TotalURLs,
		TotalClicks:     summary.TotalClicks,
		ClicksToday:     summary.ClicksToday,
		ClicksYesterday: summary.ClicksYesterday,
		AverageClicks:   summary.AverageClicks,
		TrendDirection:  trendDirection,
	}
}

func ToTopURLResponses(urls []model.URLClickStats) []dto.TopURLResponse {
	var responses []dto.TopURLResponse
	for _, url := range urls {
		responses = append(responses, dto.TopURLResponse{
			URLID:       url.URLID,
			ShortCode:   url.ShortCode,
			OriginalURL: url.OriginalURL,
			ClickCount:  url.ClickCount,
			CreatedAt:   url.CreatedAt,
		})
	}
	return responses
}

func ToReferrerResponses(referrers []model.ReferrerStats) []dto.ReferrerResponse {
	var responses []dto.ReferrerResponse
	for _, referrer := range referrers {
		responses = append(responses, dto.ReferrerResponse{
			Referrer: referrer.Referrer,
			Clicks:   referrer.Clicks,
		})
	}
	return responses
}

func ToDeviceResponses(devices []model.DeviceStats) []dto.DeviceResponse {
	var responses []dto.DeviceResponse

	var totalClicks int64
	for _, device := range devices {
		totalClicks += device.Clicks
	}

	for _, device := range devices {
		percentage := 0.0
		if totalClicks > 0 {
			percentage = float64(device.Clicks) / float64(totalClicks) * 100
		}

		responses = append(responses, dto.DeviceResponse{
			DeviceType: device.DeviceType,
			Clicks:     device.Clicks,
			Percentage: percentage,
		})
	}
	return responses
}

func ToDailyTrendResponses(trend []model.DailyClickStats) []dto.DailyTrendResponse {
	var responses []dto.DailyTrendResponse
	for _, day := range trend {
		responses = append(responses, dto.DailyTrendResponse{
			Date:   day.Date,
			Clicks: day.Clicks,
		})
	}
	return responses
}

func ToTopURLsResponse(urls []model.URLClickStats) dto.TopURLsResponse {
	return dto.TopURLsResponse{
		URLs: ToTopURLResponses(urls),
	}
}

func ToTopReferrersResponse(referrers []model.ReferrerStats) dto.TopReferrersResponse {
	return dto.TopReferrersResponse{
		Referrers: ToReferrerResponses(referrers),
	}
}

func ToDeviceBreakdownResponse(devices []model.DeviceStats) dto.DeviceBreakdownResponse {
	return dto.DeviceBreakdownResponse{
		Devices: ToDeviceResponses(devices),
	}
}

func ToDailyTrendAnalyticsResponse(trend []model.DailyClickStats, days int) dto.DailyTrendAnalyticsResponse {
	return dto.DailyTrendAnalyticsResponse{
		Trend: ToDailyTrendResponses(trend),
		Days:  days,
	}
}
