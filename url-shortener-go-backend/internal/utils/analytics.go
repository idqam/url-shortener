package utils

import (
	"net/url"
	"strings"
	"time"
	"url-shortener-go-backend/internal/model"
)


func GetDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)
	
	switch {
	case strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone"):
		return "mobile"
	case strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad"):
		return "tablet"
	case strings.Contains(ua, "bot") || strings.Contains(ua, "crawler"):
		return "bot"
	default:
		return "desktop"
	}
}


func CleanReferrer(referrer string) string {
	if referrer == "" || referrer == "direct" {
		return "direct"
	}
	
	u, err := url.Parse(referrer)
	if err != nil {
		return "unknown"
	}
	
	
	host := u.Host
	if host == "" {
		return "unknown"
	}
	
	
	host = strings.TrimPrefix(host, "www.")
	
	return host
}


func FillMissingDays(stats []model.DailyStats, days int) []model.DailyStats {
	if len(stats) == 0 {
		stats = make([]model.DailyStats, 0, days)
	}
	
	
	statsMap := make(map[string]int64)
	for _, stat := range stats {
		statsMap[stat.Date] = stat.Clicks
	}
	
	
	result := make([]model.DailyStats, 0, days)
	endDate := time.Now()
	
	for i := 0; i < days; i++ {
		date := endDate.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		
		clicks, exists := statsMap[dateStr]
		if !exists {
			clicks = 0
		}
		
		result = append(result, model.DailyStats{
			Date:   dateStr,
			Clicks: clicks,
		})
	}
	
	
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	
	return result
}