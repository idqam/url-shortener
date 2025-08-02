package model

import (
	"time"
)

type Analytics struct {
    ID         string     `json:"id" db:"id"`
    URLID      *string    `json:"url_id" db:"url_id"`
    Country    *string    `json:"country" db:"country"`
    Referrer   *string    `json:"referrer" db:"referrer"`
    UserAgent  *string    `json:"user_agent" db:"user_agent"`
    DeviceType *string    `json:"device_type" db:"device_type"`
    Browser    *string    `json:"browser" db:"browser"`
    OS         *string    `json:"os" db:"os"`
    ClickedAt  time.Time  `json:"clicked_at" db:"clicked_at"`
}

type AnalyticsSummary struct {
    TotalClicks    int64                    `json:"total_clicks"`
    UniqueClicks   int64                    `json:"unique_clicks"`
    ClicksByDay    []DailyStats             `json:"clicks_by_day"`
    TopCountries   []CountryStats           `json:"top_countries"`
    TopReferrers   []ReferrerStats          `json:"top_referrers"`
    DeviceTypes    []DeviceStats            `json:"device_types"`
    Browsers       []BrowserStats           `json:"browsers"`
}

type DailyStats struct {
    Date   string `json:"date"`
    Clicks int64  `json:"clicks"`
}

type CountryStats struct {
    Country string `json:"country"`
    Clicks  int64  `json:"clicks"`
}

type ReferrerStats struct {
    Referrer string `json:"referrer"`
    Clicks   int64  `json:"clicks"`
}

type DeviceStats struct {
    DeviceType string `json:"device_type"`
    Clicks     int64  `json:"clicks"`
}

type BrowserStats struct {
    Browser string `json:"browser"`
    Clicks  int64  `json:"clicks"`
}