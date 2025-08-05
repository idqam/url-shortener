package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"url-shortener-go-backend/internal/handler/mapper"
	"url-shortener-go-backend/internal/middleware"
	"url-shortener-go-backend/internal/service"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

func (h *AnalyticsHandler) HandleGetDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		log.Printf("[HandleGetDashboard] Fetching analytics dashboard for user: %s", userID)

		summary, err := h.analyticsService.GetUserDashboard(r.Context(), userID)
		if err != nil {
			log.Printf("[HandleGetDashboard] Failed to get dashboard for user %s: %v", userID, err)
			respondError(w, http.StatusInternalServerError, "Failed to fetch analytics dashboard")
			return
		}

		response := mapper.ToAnalyticsDashboardResponse(*summary)

		log.Printf("[HandleGetDashboard] Successfully fetched dashboard for user: %s", userID)
		RespondJSON(w, http.StatusOK, response)
	}
}

func (h *AnalyticsHandler) HandleGetTopURLs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		limit := 10
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}

		log.Printf("[HandleGetTopURLs] Fetching top %d URLs for user: %s", limit, userID)

		urls, err := h.analyticsService.GetUserTopURLs(r.Context(), userID, limit)
		if err != nil {
			log.Printf("[HandleGetTopURLs] Failed to get top URLs for user %s: %v", userID, err)
			respondError(w, http.StatusInternalServerError, "Failed to fetch top URLs")
			return
		}

		response := mapper.ToTopURLsResponse(urls)

		RespondJSON(w, http.StatusOK, response)
	}
}

func (h *AnalyticsHandler) HandleGetTopReferrers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		limit := 5
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
				limit = parsedLimit
			}
		}

		log.Printf("[HandleGetTopReferrers] Fetching top %d referrers for user: %s", limit, userID)

		referrers, err := h.analyticsService.GetUserTopReferrers(r.Context(), userID, limit)
		if err != nil {
			log.Printf("[HandleGetTopReferrers] Failed to get top referrers for user %s: %v", userID, err)
			respondError(w, http.StatusInternalServerError, "Failed to fetch top referrers")
			return
		}

		response := mapper.ToTopReferrersResponse(referrers)

		RespondJSON(w, http.StatusOK, response)
	}
}

func (h *AnalyticsHandler) HandleGetDeviceBreakdown() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		log.Printf("[HandleGetDeviceBreakdown] Fetching device breakdown for user: %s", userID)

		devices, err := h.analyticsService.GetUserDeviceBreakdown(r.Context(), userID)
		if err != nil {
			log.Printf("[HandleGetDeviceBreakdown] Failed to get device breakdown for user %s: %v", userID, err)
			respondError(w, http.StatusInternalServerError, "Failed to fetch device breakdown")
			return
		}

		response := mapper.ToDeviceBreakdownResponse(devices)

		RespondJSON(w, http.StatusOK, response)
	}
}

func (h *AnalyticsHandler) HandleGetDailyTrend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		days := 7
		if daysStr := r.URL.Query().Get("days"); daysStr != "" {
			if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 && parsedDays <= 30 {
				days = parsedDays
			}
		}

		log.Printf("[HandleGetDailyTrend] Fetching %d-day trend for user: %s", days, userID)

		trend, err := h.analyticsService.GetUserDailyTrend(r.Context(), userID, days)
		if err != nil {
			log.Printf("[HandleGetDailyTrend] Failed to get daily trend for user %s: %v", userID, err)
			respondError(w, http.StatusInternalServerError, "Failed to fetch daily trend")
			return
		}

		response := mapper.ToDailyTrendAnalyticsResponse(trend, days)

		RespondJSON(w, http.StatusOK, response)
	}
}

func (h *AnalyticsHandler) HandleRecordAnalytics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var req struct {
			URLID      string `json:"url_id"`
			Referrer   string `json:"referrer,omitempty"`
			DeviceType string `json:"device_type,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.URLID == "" {
			respondError(w, http.StatusBadRequest, "url_id is required")
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())

		log.Printf("[HandleRecordAnalytics] Recording analytics: userID=%s, urlID=%s", userID, req.URLID)

		err := h.analyticsService.RecordAnalytics(r.Context(), userID, req.URLID, req.Referrer, req.DeviceType)
		if err != nil {
			log.Printf("[HandleRecordAnalytics] Failed to record analytics: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to record analytics")
			return
		}

		RespondJSON(w, http.StatusAccepted, map[string]string{"status": "recorded"})
	}
}
