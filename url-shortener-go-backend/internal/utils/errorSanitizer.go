package utils

import (
	"strings"
)

func SanitizeError(err error, userMessage string) string {
	if err == nil {
		return userMessage
	}

	errStr := strings.ToLower(err.Error())

	sensitiveKeywords := []string{
		"supabase",
		"postgresql",
		"database",
		"relation",
		"column",
		"constraint",
		"duplicate key",
		"foreign key",
		"syntax error",
		"connection",
		"timeout",
		"internal",
	}

	for _, keyword := range sensitiveKeywords {
		if strings.Contains(errStr, keyword) {
			return userMessage
		}
	}

	return userMessage
}
