package model

import (
	"encoding/json"
	"time"
)

// DailyAnalytics represents aggregated daily statistics for a URL
type DailyAnalytics struct {
	ID              string    `json:"id" db:"id"`
	URLID           string    `json:"url_id" db:"url_id"`
	UserID          *string   `json:"user_id" db:"user_id"`
	Date            time.Time `json:"date" db:"date"`
	ClickCount      int       `json:"click_count" db:"click_count"`
	UniqueReferrers int       `json:"unique_referrers" db:"unique_referrers"`
	DesktopClicks   int       `json:"desktop_clicks" db:"desktop_clicks"`
	MobileClicks    int       `json:"mobile_clicks" db:"mobile_clicks"`
	TabletClicks    int       `json:"tablet_clicks" db:"tablet_clicks"`
	UnknownClicks   int       `json:"unknown_clicks" db:"unknown_clicks"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// UserAnalyticsSummary provides comprehensive analytics overview for user dashboard
type UserAnalyticsSummary struct {
	TotalURLs       int64             `json:"total_urls"`
	TotalClicks     int64             `json:"total_clicks"`
	ClicksToday     int64             `json:"clicks_today"`
	ClicksYesterday int64             `json:"clicks_yesterday"`
	AverageClicks   float64           `json:"average_clicks"`
	TopURLs         []URLClickStats   `json:"top_urls"`
	TopReferrers    []ReferrerStats   `json:"top_referrers"`
	DeviceBreakdown []DeviceStats     `json:"device_breakdown"`
	DailyClickTrend []DailyClickStats `json:"daily_click_trend"`
}

type URLClickStats struct {
	URLID       string `json:"url_id"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	ClickCount  int64  `json:"click_count"`
	CreatedAt   string `json:"created_at"`
}

func (u *URLClickStats) UnmarshalJSON(data []byte) error {

	var temp struct {
		ID          string `json:"id"`
		ShortCode   string `json:"short_code"`
		OriginalURL string `json:"original_url"`
		ClickCount  int64  `json:"click_count"`
		CreatedAt   string `json:"created_at"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	u.URLID = temp.ID
	u.ShortCode = temp.ShortCode
	u.OriginalURL = temp.OriginalURL
	u.ClickCount = temp.ClickCount
	u.CreatedAt = temp.CreatedAt

	return nil
}

type DailyClickStats struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

type ReferrerStats struct {
	Referrer string `json:"referrer"`
	Clicks   int64  `json:"clicks"`
}

type DeviceStats struct {
	DeviceType string `json:"device_type"`
	Clicks     int64  `json:"clicks"`
}

type AnalyticsDateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
