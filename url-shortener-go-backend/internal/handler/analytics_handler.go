package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"url-shortener-go-backend/internal/handler/dto"
	"url-shortener-go-backend/internal/service"
	"url-shortener-go-backend/internal/utils"

	"github.com/go-chi/chi/v5"
)

type AnalyticsHandler struct {
	service service.AnalyticsService
}

func NewAnalyticsHandler(svc service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: svc}
}

// GET /analytics/url/{urlID}?days=30
func (h *AnalyticsHandler) GetURLAnalytics(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "urlID")
	days := getQueryInt(r, "days", 30)

	analytics, err := h.service.GetURLAnalytics(r.Context(), urlID, days)
	if err != nil {
		http.Error(w, "Failed to get analytics", http.StatusInternalServerError)
		return
	}

	resp := dto.NewAnalyticsSummaryResponse(analytics)
	writeJSON(w, http.StatusOK, resp)
}

// GET /analytics/user/summary?days=30
func (h *AnalyticsHandler) GetUserAnalytics(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	days := getQueryInt(r, "days", 30)
	analytics, err := h.service.GetUserAnalytics(r.Context(), userID, days)
	if err != nil {
		http.Error(w, "Failed to get analytics", http.StatusInternalServerError)
		return
	}

	resp := dto.NewAnalyticsSummaryResponse(analytics)
	writeJSON(w, http.StatusOK, resp)
}

// GET /analytics/user/top-urls?limit=10
func (h *AnalyticsHandler) GetTopURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit := getQueryInt(r, "limit", 10)
	stats, err := h.service.GetTopURLs(r.Context(), userID, limit)
	if err != nil {
		http.Error(w, "Failed to get top URLs", http.StatusInternalServerError)
		return
	}
	resp := dto.NewURLStatsResponse(stats)
writeJSON(w, http.StatusOK, resp)

}

func (h *AnalyticsHandler) GetUserDailyStats(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	days := getQueryInt(r, "days", 30)
	urlID := r.URL.Query().Get("url_id")

	var stats interface{}

	if urlID != "" {
		stats, err = h.service.GetDailyStats(r.Context(), &urlID, nil, days)
	} else {
		stats, err = h.service.GetDailyStats(r.Context(), nil, &userID, days)
	}

	if err != nil {
		http.Error(w, "Failed to get daily stats", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (h *AnalyticsHandler) GetUserDeviceStats(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urlID := r.URL.Query().Get("url_id")

	var stats interface{}

	if urlID != "" {
		stats, err = h.service.GetDeviceStats(r.Context(), &urlID, nil)
	} else {
		stats, err = h.service.GetDeviceStats(r.Context(), nil, &userID)
	}

	if err != nil {
		http.Error(w, "Failed to get device stats", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (h *AnalyticsHandler) GetUserReferrerStats(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit := getQueryInt(r, "limit", 10)
	urlID := r.URL.Query().Get("url_id")

	var stats interface{}

	if urlID != "" {
		stats, err = h.service.GetReferrerStats(r.Context(), &urlID, nil, limit)
	} else {
		stats, err = h.service.GetReferrerStats(r.Context(), nil, &userID, limit)
	}

	if err != nil {
		http.Error(w, "Failed to get referrer stats", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}



func extractUserIDFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization format")
	}
	return utils.ExtractUserIDFromSupabaseToken(parts[1])
}

func getQueryInt(r *http.Request, key string, defaultVal int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultVal
	}
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}
