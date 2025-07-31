package utils

func IsValidShortCode(code string) bool {
	if len(code) < 6 || len(code) > 12 {
		return false
	}
	for _, c := range code {
		if !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') && !(c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}
