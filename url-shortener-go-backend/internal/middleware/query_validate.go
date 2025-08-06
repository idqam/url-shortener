package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"url-shortener-go-backend/internal/utils"
)

func ValidateQueryParams() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limit := r.URL.Query().Get("limit"); limit != "" {
				if val, err := strconv.Atoi(limit); err != nil || val < 1 || val > 100 {
					http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
					return
				}
			}

			if days := r.URL.Query().Get("days"); days != "" {
				if val, err := strconv.Atoi(days); err != nil || val < 1 || val > 365 {
					http.Error(w, "Invalid days parameter", http.StatusBadRequest)
					return
				}
			}

			if strings.HasPrefix(r.URL.Path, "/api/urls/") {
				shortcode := strings.TrimPrefix(r.URL.Path, "/api/urls/")
				if !utils.IsValidShortCode(shortcode) {
					http.Error(w, "Invalid shortcode format", http.StatusBadRequest)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
