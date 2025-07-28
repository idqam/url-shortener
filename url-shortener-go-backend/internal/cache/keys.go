package cache

import "fmt"

func KeyShortCode(sc string) string {
	return fmt.Sprintf("short:%s", sc)
}

func KeyOriginalURL(url string) string {
	return fmt.Sprintf("url:%s", url)
}
