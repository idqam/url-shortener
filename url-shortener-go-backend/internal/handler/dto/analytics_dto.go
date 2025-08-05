package dto

type AnalyticsDashboardResponse struct {
	Overview        AnalyticsOverview    `json:"overview"`
	TopURLs         []TopURLResponse     `json:"top_urls"`
	TopReferrers    []ReferrerResponse   `json:"top_referrers"`
	DeviceBreakdown []DeviceResponse     `json:"device_breakdown"`
	DailyTrend      []DailyTrendResponse `json:"daily_trend"`
}

type AnalyticsOverview struct {
	TotalURLs       int64   `json:"total_urls"`
	TotalClicks     int64   `json:"total_clicks"`
	ClicksToday     int64   `json:"clicks_today"`
	ClicksYesterday int64   `json:"clicks_yesterday"`
	AverageClicks   float64 `json:"average_clicks"`
	TrendDirection  string  `json:"trend_direction"`
}

type TopURLResponse struct {
	URLID       string `json:"url_id"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	ClickCount  int64  `json:"click_count"`
	CreatedAt   string `json:"created_at"`
}

type ReferrerResponse struct {
	Referrer string `json:"referrer"`
	Clicks   int64  `json:"clicks"`
}

type DeviceResponse struct {
	DeviceType string  `json:"device_type"`
	Clicks     int64   `json:"clicks"`
	Percentage float64 `json:"percentage"`
}

type DailyTrendResponse struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

type AnalyticsRequest struct {
	Days  int `json:"days,omitempty"`
	Limit int `json:"limit,omitempty"`
}

type TopURLsResponse struct {
	URLs []TopURLResponse `json:"urls"`
}

type TopReferrersResponse struct {
	Referrers []ReferrerResponse `json:"referrers"`
}

type DeviceBreakdownResponse struct {
	Devices []DeviceResponse `json:"devices"`
}

type DailyTrendAnalyticsResponse struct {
	Trend []DailyTrendResponse `json:"trend"`
	Days  int                  `json:"days"`
}
