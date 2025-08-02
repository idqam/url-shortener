// File: handler/url_handler.go
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
	"url-shortener-go-backend/internal/middleware"
	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/service"
	"url-shortener-go-backend/internal/utils"
)

type URLHandler struct {
	svc   service.URLService
	cache cache.Cache
}

func NewURLHandler(svc service.URLService, cache cache.Cache) *URLHandler {
	return &URLHandler{
		svc:   svc,
		cache: cache,
	}
}

// GET /<shortcode>
func (s *URLHandler) ShortCodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortcode := strings.Trim(r.URL.Path, "/")
		if utils.IsValidShortCode(shortcode) {
			s.HandleRedirect()(w, r)
			return
		}
		http.NotFound(w, r)
	}
}

// POST /api/shorten
func (h *URLHandler) HandleShorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req dto.ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		req.URL = strings.TrimSpace(req.URL)
		if req.URL == "" {
			respondError(w, http.StatusBadRequest, "URL is required")
			return
		}

		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		} else {
		
			ip = strings.Split(ip, ",")[0]
		}

		authHeader := r.Header.Get("Authorization")
		var userID *string
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := middleware.ExtractClaimsFromToken(token)
			if err == nil {
				if sub, ok := claims["sub"].(string); ok && sub != "" {
					userID = &sub
				}
			}
		}

		req.UserID = userID

		userIDPart := "anon"
		if userID != nil && *userID != "" {
			userIDPart = *userID
		}
		rateKey := fmt.Sprintf("rate_limit:%s:%s", userIDPart, ip)

		ctx := r.Context()
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
			log.Printf("[RateLimit] Too many requests from %s (%s): %d", userIDPart, ip, count)
			http.Error(w, "Too many requests — slow down", http.StatusTooManyRequests)
			return
		}

		result, err := h.svc.CreateShortURL(ctx, req.URL, req.UserID, req.IsPublic, req.CodeLength)
		if err != nil {
			log.Printf("[HandleShorten] CreateShortURL failed: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to create short URL")
			return
		}

		res := dto.ToShortenResponse(result)
		respondJSON(w, http.StatusCreated, res)
	}
}

// GET /api/urls
func (h *URLHandler) HandleGetUserUrls() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			respondError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := utils.ExtractUserIDFromSupabaseToken(token)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		user := model.User{ID: userID}
		urls, err := h.svc.GetUserUrls(user)
		if err != nil {
			log.Printf("[HandleGetUserUrls] GetUserUrls failed: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to fetch URLs")
			return
		}

		response := dto.ToGetUrlsResponse(urls)
		respondJSON(w, http.StatusOK, response)
	}
}

// GET /<shortcode>
func (h *URLHandler) HandleRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		shortcode := strings.Trim(r.URL.Path, "/")
		if shortcode == "" {
			respondError(w, http.StatusBadRequest, "Shortcode is required")
			return
		}

		ctx := r.Context()
		urlEntry, err := h.svc.GetURLByShortCode(ctx, shortcode)
		if err != nil || urlEntry == nil || urlEntry.OriginalURL == "" {
			log.Printf("[HandleRedirect] URL not found for shortcode: %s", shortcode)
			respondError(w, http.StatusNotFound, "URL not found")
			return
		}

		go func() {
			if err := h.svc.IncrementClickCount(shortcode); err != nil {
				log.Printf("[HandleRedirect] Failed to increment click count: %v", err)
			}
		}()

		http.Redirect(w, r, urlEntry.OriginalURL, http.StatusFound)
	}
}

// GET /api/url?shortcode=abc123
func (h *URLHandler) HandleGetUrlByShortCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		shortcode := strings.TrimSpace(r.URL.Query().Get("shortcode"))
		if shortcode == "" {
			respondError(w, http.StatusBadRequest, "Shortcode is required")
			return
		}

		ctx := r.Context()
		u, err := h.svc.GetURLByShortCode(ctx, shortcode)
		if err != nil || u == nil {
			log.Printf("[HandleGetUrlByShortCode] Not found: %s", shortcode)
			respondError(w, http.StatusNotFound, "URL not found")
			return
		}

		respondJSON(w, http.StatusOK, dto.ToURLResponse(*u))
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	log.Printf("[respondError] %d - %s", status, message)
	respondJSON(w, status, dto.ErrorResponse{Error: message})
}
