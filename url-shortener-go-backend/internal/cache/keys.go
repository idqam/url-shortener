package cache

import (
	"fmt"
	"time"
)

func KeyShortCode(sc string) string {
	return fmt.Sprintf("short:%s", sc)
}

func KeyOriginalURL(url string) string {
	return fmt.Sprintf("url:%s", url)
}

func KeyURLAnalytics(urlID string, startDate, endDate time.Time) string {
	return fmt.Sprintf("analytics:summary:url:%s:%s:%s",
		urlID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
}

func KeyUserAnalytics(userID string, startDate, endDate time.Time) string {
	return fmt.Sprintf("analytics:summary:user:%s:%s:%s",
		userID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
}
