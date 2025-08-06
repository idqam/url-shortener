package cache

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"
)

func SecureKey(keyType, identifier string, additionalData ...string) string {
	salt := os.Getenv("SALT")
	if salt == "" {
		salt = "default-cache-salt-change-in-production"
	}

	data := fmt.Sprintf("%s:%s:%s", keyType, identifier, salt)
	for _, extra := range additionalData {
		data += ":" + extra
	}

	hasher := sha256.New()
	hasher.Write([]byte(data))
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	return fmt.Sprintf("%s_%s", keyType, hash[:16])
}

func KeyShortCode(sc string) string {
	return SecureKey("short", sc)
}

func KeyOriginalURL(url string) string {
	return SecureKey("url", url)
}

func KeyURLAnalytics(urlID string, startDate, endDate time.Time) string {
	return SecureKey("analytics:summary:url:%s:%s:%s",
		urlID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
}

func KeyUserAnalytics(userID string, startDate, endDate time.Time) string {
	return SecureKey("analytics", userID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
}
