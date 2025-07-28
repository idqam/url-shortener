package handler

import (
	"encoding/json"
	"fmt"
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
func (s *URLHandler) ShortCodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := strings.Trim(r.URL.Path, "/")
		if utils.IsValidShortCode(code) {
			s.HandleRedirect()(w, r)
			return
		}
		http.NotFound(w, r)
	}
}

func (h *URLHandler) HandleShorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		result, err := h.svc.CreateShortURL(r.Context(), req.URL, req.UserID, req.IsPublic, req.CodeLength)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		res := dto.ToShortenResponse(result)
		respondJSON(w, http.StatusCreated, res)
	}
}

//GET /api/urls?user_id=123
func (h *URLHandler) HandleGetUserUrls() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			respondError(w, http.StatusBadRequest, "missing user_id parameter")
			return
		}

		user := model.User{ID: userID}
		urls, err := h.svc.GetUserUrls(user)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to fetch URLs")
			return
		}

		response := dto.ToGetUrlsResponse(urls)
		respondJSON(w, http.StatusOK, response)
	}
}

//GET /<shortcode>
func (h *URLHandler) HandleRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortcode := strings.Trim(r.URL.Path, "/")
		if shortcode == "" {
			respondError(w, http.StatusBadRequest, "Shortcode is required")
			return
		}

		ctx := r.Context()
		urlEntry, err := h.svc.GetURLByShortCode(ctx, shortcode)
		if err != nil || urlEntry == nil || urlEntry.OriginalURL == "" {
			respondError(w, http.StatusNotFound, "URL not found")
			return
		}

		go func() {
			if err := h.svc.IncrementClickCount(shortcode); err != nil {
				fmt.Printf("failed to increment click count: %v\n", err)
			}
		}()

		http.Redirect(w, r, urlEntry.OriginalURL, http.StatusFound)
	}
}

// GET /api/url?shortcode=abc123
func (h *URLHandler) HandleGetUrlByShortCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortcode := r.URL.Query().Get("shortcode")
		if shortcode == "" {
			respondError(w, http.StatusBadRequest, "shortcode is required")
			return
		}

		ctx := r.Context()
		u, err := h.svc.GetURLByShortCode(ctx, shortcode)
		if err != nil || u == nil {
			respondError(w, http.StatusNotFound, "URL not found")
			return
		}

		respondJSON(w, http.StatusOK, dto.ToURLResponse(*u))
	}
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, dto.ErrorResponse{Error: message})
}
