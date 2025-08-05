package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"url-shortener-go-backend/internal/cache"
	"url-shortener-go-backend/internal/handler/dto"
	"url-shortener-go-backend/internal/handler/mapper"
	"url-shortener-go-backend/internal/middleware"
	"url-shortener-go-backend/internal/service"
	"url-shortener-go-backend/internal/utils"
)

type URLHandler struct {
	svc   service.URLService
	cache cache.Cache
}

func NewURLHandler(svc service.URLService, cache cache.Cache) *URLHandler {
	return &URLHandler{svc: svc, cache: cache}
}

func (h *URLHandler) HandleShorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req dto.ShortenURLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.OriginalURL) == "" {
			respondError(w, http.StatusBadRequest, "Invalid or missing URL")
			return
		}

		ctx := r.Context()
		userID := middleware.GetUserIDFromContext(ctx)
		var userIDPtr *string
		if userID != "" {
			userIDPtr = &userID
		}

		ip := getClientIP(r)
		rateKey := fmt.Sprintf("rate_limit:%s:%s", getUserIDOrAnonymous(userID), ip)
		count, err := h.cache.Incr(ctx, rateKey)
		if err != nil {
			log.Printf("[RateLimit] INCR failed: %v", err)
			http.Error(w, "Rate limit error", http.StatusInternalServerError)
			return
		}
		if count == 1 {
			_ = h.cache.Expire(ctx, rateKey, 60*time.Second)
		}
		if count > 10 {
			http.Error(w, "Too many requests — slow down", http.StatusTooManyRequests)
			return
		}

		urlModel, err := h.svc.CreateShortURL(ctx, strings.TrimSpace(req.OriginalURL), req.IsPublic, userIDPtr, int(req.CodeLength))
		if err != nil {
			log.Printf("[CreateShortURL] failed: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to shorten URL")
			return
		}

		RespondJSON(w, http.StatusCreated, mapper.ToShortenURLResponse(*urlModel))
	}
}

func (h *URLHandler) HandleGetUserUrls() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		urls, err := h.svc.GetUserUrls(r.Context(), userID)
		if err != nil {
			log.Printf("[GetUserUrls] failed: %v", err)
			respondError(w, http.StatusInternalServerError, "Could not fetch URLs")
			return
		}

		RespondJSON(w, http.StatusOK, mapper.ToGetUrlsResponse(urls))
	}
}

func (h *URLHandler) HandleGetUrlByShortCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		shortcode := strings.TrimPrefix(r.URL.Path, "/api/urls/")
		shortcode = strings.TrimSpace(shortcode)
		if shortcode == "" {
			respondError(w, http.StatusBadRequest, "Shortcode is required")
			return
		}

		ctx := r.Context()
		url, err := h.svc.GetURLByShortCode(ctx, shortcode)
		if err != nil || url == nil {
			respondError(w, http.StatusNotFound, "URL not found")
			return
		}

		RespondJSON(w, http.StatusOK, mapper.ToGetURLByShortCodeResponse(*url))
	}
}

func (h *URLHandler) ShortCodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortcode := strings.Trim(r.URL.Path, "/")
		if utils.IsValidShortCode(shortcode) {
			h.HandleRedirect()(w, r)
			return
		}
		http.NotFound(w, r)
	}
}

func (h *URLHandler) HandleRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		shortcode := strings.Trim(r.URL.Path, "/")
		if shortcode == "" {
			respondError(w, http.StatusBadRequest, "Missing shortcode")
			return
		}

		ctx := r.Context()
		urlEntry, err := h.svc.GetURLByShortCode(ctx, shortcode)
		if err != nil || urlEntry == nil || urlEntry.OriginalURL == "" {
			respondError(w, http.StatusNotFound, "URL not found")
			return
		}

		go func() {

			if err := h.svc.IncrementClickCount(ctx, shortcode); err != nil {
				log.Printf("[ClickCount] Increment failed: %v", err)
			}
		}()

		http.Redirect(w, r, urlEntry.OriginalURL, http.StatusFound)
	}
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	} else {
		ip = strings.Split(ip, ",")[0]
	}
	return strings.TrimSpace(ip)
}

func getUserIDOrAnonymous(userID string) string {
	if userID == "" {
		return "anonymous"
	}
	return userID
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	log.Printf("[Error] %d: %s", status, message)
	RespondJSON(w, status, dto.ErrorResponse{Error: message})
}
