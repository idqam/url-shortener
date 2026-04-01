package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
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
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

func (h *AnalyticsHandler) HandleGetDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetRequestID(r.Context())

		if r.Method != http.MethodGet {
			h.respondError(w, r, http.StatusMethodNotAllowed, ErrMsgMethodNotAllowed, requestID)
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			h.respondError(w, r, http.StatusUnauthorized, ErrMsgUnauthorized, requestID)
			return
		}

		slog.Info("fetching dashboard", "request_id", requestID, "user_id", truncateID(userID))

		summary, err := h.analyticsService.GetUserDashboard(r.Context(), userID)
		if err != nil {
			slog.Error("dashboard fetch failed", "request_id", requestID, "user_id", truncateID(userID), "error", err)
			h.respondError(w, r, http.StatusInternalServerError,
				utils.SanitizeError(err, ErrMsgDashboardFetch), requestID)
			return
		}

		if summary == nil {
			slog.Warn("dashboard returned nil", "request_id", requestID, "user_id", truncateID(userID))
			h.respondError(w, r, http.StatusInternalServerError, ErrMsgDashboardFetch, requestID)
			return
		}

		response := mapper.ToAnalyticsDashboardResponse(*summary)

		h.respondJSON(w, http.StatusOK, response, requestID)
	}
}

func (h *AnalyticsHandler) HandleGetTopURLs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetRequestID(r.Context())

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

		slog.Info("fetching top urls", "request_id", requestID, "user_id", truncateID(userID), "limit", limit)

		urls, err := h.analyticsService.GetUserTopURLs(r.Context(), userID, limit)
		if err != nil {
			slog.Error("top urls fetch failed", "request_id", requestID, "user_id", truncateID(userID), "error", err)
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
		requestID := middleware.GetRequestID(r.Context())

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

		slog.Info("fetching top referrers", "request_id", requestID, "user_id", truncateID(userID), "limit", limit)

		referrers, err := h.analyticsService.GetUserTopReferrers(r.Context(), userID, limit)
		if err != nil {
			slog.Error("top referrers fetch failed", "request_id", requestID, "user_id", truncateID(userID), "error", err)
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
		requestID := middleware.GetRequestID(r.Context())

		if r.Method != http.MethodGet {
			h.respondError(w, r, http.StatusMethodNotAllowed, ErrMsgMethodNotAllowed, requestID)
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			h.respondError(w, r, http.StatusUnauthorized, ErrMsgUnauthorized, requestID)
			return
		}

		slog.Info("fetching device breakdown", "request_id", requestID, "user_id", truncateID(userID))

		devices, err := h.analyticsService.GetUserDeviceBreakdown(r.Context(), userID)
		if err != nil {
			slog.Error("device breakdown fetch failed", "request_id", requestID, "user_id", truncateID(userID), "error", err)
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
		requestID := middleware.GetRequestID(r.Context())

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

		slog.Info("fetching daily trend", "request_id", requestID, "user_id", truncateID(userID), "days", days)

		trend, err := h.analyticsService.GetUserDailyTrend(r.Context(), userID, days)

		if err != nil {
			slog.Error("daily trend fetch failed", "request_id", requestID, "user_id", truncateID(userID), "error", err)
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
		requestID := middleware.GetRequestID(r.Context())

		if r.Method != http.MethodPost {
			h.respondError(w, r, http.StatusMethodNotAllowed, ErrMsgMethodNotAllowed, requestID)
			return
		}

		var req dto.RecordAnalyticsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error("failed to decode analytics request", "request_id", requestID, "error", err)
			h.respondError(w, r, http.StatusBadRequest, ErrMsgInvalidRequest, requestID)
			return
		}

		if err := h.validateAnalyticsRequest(&req); err != nil {
			slog.Warn("invalid analytics request", "request_id", requestID, "error", err)
			h.respondError(w, r, http.StatusBadRequest, err.Error(), requestID)
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())

		slog.Info("recording analytics", "request_id", requestID, "user_id", truncateID(userID), "url_id", req.URLID, "device", req.DeviceType)

		err := h.analyticsService.RecordAnalytics(
			r.Context(),
			userID,
			req.URLID,
			req.Referrer,
			req.DeviceType,
		)

		if err != nil {
			slog.Error("failed to record analytics", "request_id", requestID, "error", err)
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
		slog.Error("failed to encode json response", "request_id", requestID, "error", err)
	}
}

func (h *AnalyticsHandler) respondError(w http.ResponseWriter, r *http.Request, status int, message string, requestID string) {
	slog.Info("response error", "request_id", requestID, "status", status, "path", r.URL.Path, "method", r.Method)

	h.respondJSON(w, status, dto.ErrorResponse{
		Error:     message,
		RequestID: requestID,
		Timestamp: time.Now().Unix(),
	}, requestID)
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
