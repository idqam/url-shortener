package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"url-shortener-go-backend/internal/handler/dto"
	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/service"
	"url-shortener-go-backend/internal/utils"
)

type URLHandler struct {
	svc service.URLService
}

func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
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

		result, err := h.svc.CreateShortURL(r.Context(), req.URL, req.UserID, req.IsPublic, req.CodeLength)
		if err != nil {
			log.Printf("[HandleShorten] CreateShortURL failed: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to create short URL")
			return
		}

		res := dto.ToShortenResponse(result)
		respondJSON(w, http.StatusCreated, res)
	}
}

// GET /api/urls?user_id=123
func (h *URLHandler) HandleGetUserUrls() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
		if userID == "" {
			respondError(w, http.StatusBadRequest, "Missing user_id parameter")
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
	// CORS for dev/testing
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
