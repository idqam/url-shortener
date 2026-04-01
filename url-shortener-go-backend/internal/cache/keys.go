package cache

import (
	"crypto/sha256"
	"fmt"
	"time"
)

func SecureKey(salt, keyType, identifier string, additionalData ...string) string {
	data := fmt.Sprintf("%s:%s:%s", keyType, identifier, salt)
	for _, extra := range additionalData {
		data += ":" + extra
	}

	hasher := sha256.New()
	hasher.Write([]byte(data))
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	return fmt.Sprintf("%s_%s", keyType, hash[:16])
}

func KeyShortCode(salt, sc string) string {
	return SecureKey(salt, "short", sc)
}

func KeyOriginalURL(salt, url string) string {
	return SecureKey(salt, "url", url)
}

func KeyURLAnalytics(salt, urlID string, startDate, endDate time.Time) string {
	return SecureKey(salt, "analytics:summary:url:%s:%s:%s",
		urlID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
}

func KeyUserAnalytics(salt, userID string, startDate, endDate time.Time) string {
	return SecureKey(salt, "analytics", userID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
}
