package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"url-shortener-go-backend/internal/handler/dto"
	"url-shortener-go-backend/internal/handler/mapper"
	"url-shortener-go-backend/internal/middleware"
	"url-shortener-go-backend/internal/service"
	"url-shortener-go-backend/internal/utils"
)

const (
	ErrMsgMethodNotAllowed = "Method not allowed"
	ErrMsgUnauthorized     = "Authentication required"
	ErrMsgInvalidLimit     = "Invalid limit parameter"
	ErrMsgInvalidDays      = "Invalid days parameter"
	ErrMsgInvalidRequest   = "Invalid request format"
	ErrMsgInvalidURLID     = "Invalid URL identifier"
	ErrMsgInternalError    = "An error occurred while processing your request"
	ErrMsgDashboardFetch   = "Unable to fetch analytics dashboard"
	ErrMsgTopURLsFetch     = "Unable to fetch top URLs"
	ErrMsgReferrersFetch   = "Unable to fetch referrer data"
	ErrMsgDevicesFetch     = "Unable to fetch device breakdown"
	ErrMsgTrendFetch       = "Unable to fetch trend data"
	ErrMsgRecordFailed     = "Unable to record analytics data"

	MaxLimit              = 100
	MinLimit              = 1
	MaxDays               = 365
	MinDays               = 1
	MaxReferrerLength     = 500
	MaxURLIDLength        = 50
	DefaultTopURLsLimit   = 10
	DefaultReferrersLimit = 5
	DefaultTrendDays      = 7
)

var (
	validDeviceTypes = map[string]bool{
		"desktop": true,
		"mobile":  true,
		"tablet":  true,
		"bot":     true,
		"unknown": true,
	}
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
	requestIDGen     func() string
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		requestIDGen:     generateRequestID,
	}
}

func (h *AnalyticsHandler) HandleGetDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := h.requestIDGen()

		if r.Method != http.MethodGet {
			h.respondError(w, r, http.StatusMethodNotAllowed, ErrMsgMethodNotAllowed, requestID)
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			h.respondError(w, r, http.StatusUnauthorized, ErrMsgUnauthorized, requestID)
			return
		}

		log.Printf("[%s] Fetching dashboard for user: %s", requestID, truncateID(userID))

		summary, err := h.analyticsService.GetUserDashboard(r.Context(), userID)
		service.NormalizeSummary(summary)
		if err != nil {
			log.Printf("[%s] ERROR: Dashboard fetch failed for user %s: %v",
				requestID, truncateID(userID), err)
			h.respondError(w, r, http.StatusInternalServerError,
				utils.SanitizeError(err, ErrMsgDashboardFetch), requestID)
			return
		}

		if summary == nil {
			log.Printf("[%s] WARN: Dashboard returned nil for user %s", requestID, truncateID(userID))
			h.respondError(w, r, http.StatusInternalServerError, ErrMsgDashboardFetch, requestID)
			return
		}

		response := mapper.ToAnalyticsDashboardResponse(*summary)
		h.respondJSON(w, http.StatusOK, response, requestID)
	}
}

func (h *AnalyticsHandler) HandleGetTopURLs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := h.requestIDGen()

		if r.Method != http.MethodGet {
			h.respondError(w, r, http.StatusMethodNotAllowed, ErrMsgMethodNotAllowed, requestID)
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			h.respondError(w, r, http.StatusUnauthorized, ErrMsgUnauthorized, requestID)
			return
		}

		limit, err := h.parseLimit(r.URL.Query().Get("limit"), DefaultTopURLsLimit, MaxLimit)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, ErrMsgInvalidLimit, requestID)
			return
		}

		log.Printf("[%s] Fetching top %d URLs for user: %s", requestID, limit, truncateID(userID))

		urls, err := h.analyticsService.GetUserTopURLs(r.Context(), userID, limit)
		if err != nil {
			log.Printf("[%s] ERROR: Top URLs fetch failed for user %s: %v",
				requestID, truncateID(userID), err)
			h.respondError(w, r, http.StatusInternalServerError,
				utils.SanitizeError(err, ErrMsgTopURLsFetch), requestID)
			return
		}

		response := mapper.ToTopURLsResponse(urls)
		h.respondJSON(w, http.StatusOK, response, requestID)
	}
}

func (h *AnalyticsHandler) HandleGetTopReferrers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := h.requestIDGen()

		if r.Method != http.MethodGet {
			h.respondError(w, r, http.StatusMethodNotAllowed, ErrMsgMethodNotAllowed, requestID)
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			h.respondError(w, r, http.StatusUnauthorized, ErrMsgUnauthorized, requestID)
			return
		}

		limit, err := h.parseLimit(r.URL.Query().Get("limit"), DefaultReferrersLimit, 50)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, ErrMsgInvalidLimit, requestID)
			return
		}

		log.Printf("[%s] Fetching top %d referrers for user: %s", requestID, limit, truncateID(userID))

		referrers, err := h.analyticsService.GetUserTopReferrers(r.Context(), userID, limit)
		if err != nil {
			log.Printf("[%s] ERROR: Top referrers fetch failed for user %s: %v",
				requestID, truncateID(userID), err)
			h.respondError(w, r, http.StatusInternalServerError,
				utils.SanitizeError(err, ErrMsgReferrersFetch), requestID)
			return
		}

		response := mapper.ToTopReferrersResponse(referrers)
		h.respondJSON(w, http.StatusOK, response, requestID)
	}
}

func (h *AnalyticsHandler) HandleGetDeviceBreakdown() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := h.requestIDGen()

		if r.Method != http.MethodGet {
			h.respondError(w, r, http.StatusMethodNotAllowed, ErrMsgMethodNotAllowed, requestID)
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			h.respondError(w, r, http.StatusUnauthorized, ErrMsgUnauthorized, requestID)
			return
		}

		log.Printf("[%s] Fetching device breakdown for user: %s", requestID, truncateID(userID))

		devices, err := h.analyticsService.GetUserDeviceBreakdown(r.Context(), userID)
		if err != nil {
			log.Printf("[%s] ERROR: Device breakdown fetch failed for user %s: %v",
				requestID, truncateID(userID), err)
			h.respondError(w, r, http.StatusInternalServerError,
				utils.SanitizeError(err, ErrMsgDevicesFetch), requestID)
			return
		}

		response := mapper.ToDeviceBreakdownResponse(devices)
		h.respondJSON(w, http.StatusOK, response, requestID)
	}
}

func (h *AnalyticsHandler) HandleGetDailyTrend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := h.requestIDGen()

		if r.Method != http.MethodGet {
			h.respondError(w, r, http.StatusMethodNotAllowed, ErrMsgMethodNotAllowed, requestID)
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			h.respondError(w, r, http.StatusUnauthorized, ErrMsgUnauthorized, requestID)
			return
		}

		days, err := h.parseDays(r.URL.Query().Get("days"), DefaultTrendDays)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, ErrMsgInvalidDays, requestID)
			return
		}

		log.Printf("[%s] Fetching %d-day trend for user: %s", requestID, days, truncateID(userID))

		trend, err := h.analyticsService.GetUserDailyTrend(r.Context(), userID, days)

		if err != nil {
			log.Printf("[%s] ERROR: Daily trend fetch failed for user %s: %v",
				requestID, truncateID(userID), err)
			h.respondError(w, r, http.StatusInternalServerError,
				utils.SanitizeError(err, ErrMsgTrendFetch), requestID)
			return
		}

		response := mapper.ToDailyTrendAnalyticsResponse(trend, days)
		h.respondJSON(w, http.StatusOK, response, requestID)
	}
}

func (h *AnalyticsHandler) HandleRecordAnalytics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := h.requestIDGen()

		if r.Method != http.MethodPost {
			h.respondError(w, r, http.StatusMethodNotAllowed, ErrMsgMethodNotAllowed, requestID)
			return
		}

		var req dto.RecordAnalyticsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[%s] ERROR: Failed to decode analytics request: %v", requestID, err)
			h.respondError(w, r, http.StatusBadRequest, ErrMsgInvalidRequest, requestID)
			return
		}

		if err := h.validateAnalyticsRequest(&req); err != nil {
			log.Printf("[%s] WARN: Invalid analytics request: %v", requestID, err)
			h.respondError(w, r, http.StatusBadRequest, err.Error(), requestID)
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())

		log.Printf("[%s] Recording analytics: user=%s, url=%s, device=%s",
			requestID, truncateID(userID), req.URLID, req.DeviceType)

		err := h.analyticsService.RecordAnalytics(
			r.Context(),
			userID,
			req.URLID,
			req.Referrer,
			req.DeviceType,
		)

		if err != nil {
			log.Printf("[%s] ERROR: Failed to record analytics: %v", requestID, err)
			h.respondError(w, r, http.StatusInternalServerError,
				utils.SanitizeError(err, ErrMsgRecordFailed), requestID)
			return
		}

		h.respondJSON(w, http.StatusAccepted, dto.RecordAnalyticsResponse{
			Status:    "recorded",
			Timestamp: time.Now().Unix(),
			RequestID: requestID,
		}, requestID)
	}
}

func (h *AnalyticsHandler) parseLimit(limitStr string, defaultLimit, maxLimit int) (int, error) {
	if limitStr == "" {
		return defaultLimit, nil
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, fmt.Errorf("must be a number")
	}

	if limit < MinLimit || limit > maxLimit {
		return 0, fmt.Errorf("must be between %d and %d", MinLimit, maxLimit)
	}

	return limit, nil
}

func (h *AnalyticsHandler) parseDays(daysStr string, defaultDays int) (int, error) {
	if daysStr == "" {
		return defaultDays, nil
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return 0, fmt.Errorf("must be a number")
	}

	if days < MinDays || days > MaxDays {
		return 0, fmt.Errorf("must be between %d and %d", MinDays, MaxDays)
	}

	return days, nil
}

func (h *AnalyticsHandler) validateAnalyticsRequest(req *dto.RecordAnalyticsRequest) error {
	req.URLID = strings.TrimSpace(req.URLID)
	if req.URLID == "" {
		return fmt.Errorf("url_id is required")
	}

	if len(req.URLID) > MaxURLIDLength {
		return fmt.Errorf("url_id is too long")
	}

	for _, char := range req.URLID {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return fmt.Errorf("url_id contains invalid characters")
		}
	}

	req.Referrer = h.sanitizeReferrer(req.Referrer)
	if req.DeviceType != "" {
		deviceType := strings.ToLower(req.DeviceType)
		if !validDeviceTypes[deviceType] {
			req.DeviceType = "unknown"
		} else {
			req.DeviceType = deviceType
		}
	} else {
		req.DeviceType = "unknown"
	}

	return nil
}

func (h *AnalyticsHandler) sanitizeReferrer(referrer string) string {
	referrer = strings.TrimSpace(referrer)

	if referrer == "" {
		return ""
	}

	dangerousSchemes := []string{"javascript:", "data:", "vbscript:", "file:"}
	lowerReferrer := strings.ToLower(referrer)
	for _, scheme := range dangerousSchemes {
		if strings.HasPrefix(lowerReferrer, scheme) {
			return ""
		}
	}

	if len(referrer) > MaxReferrerLength {
		referrer = referrer[:MaxReferrerLength]
	}

	return referrer
}

func (h *AnalyticsHandler) respondJSON(w http.ResponseWriter, status int, data interface{}, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[%s] ERROR: Failed to encode JSON response: %v", requestID, err)
	}
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func (h *AnalyticsHandler) respondError(w http.ResponseWriter, r *http.Request, status int, message string, requestID string) {
	log.Printf("[%s] Response Error: status=%d, path=%s, method=%s",
		requestID, status, r.URL.Path, r.Method)

	h.respondJSON(w, status, ErrorResponse2{
		Error:     message,
		RequestID: requestID,
		Timestamp: time.Now().Unix(),
	}, requestID)
}

func generateRequestID() string {
	return fmt.Sprintf("req_%d_%d", time.Now().UnixNano(), randomInt())
}

func truncateID(id string) string {
	if id == "" {
		return "anonymous"
	}
	if len(id) <= 8 {
		return id
	}
	return id[:8] + "***"
}

func randomInt() int {

	nBig, _ := rand.Int(rand.Reader, big.NewInt(27))
	return int(nBig.Int64())
}

type ErrorResponse2 struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
