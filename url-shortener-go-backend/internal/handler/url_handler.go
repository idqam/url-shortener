package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"url-shortener-go-backend/internal/handler/dto"
	"url-shortener-go-backend/internal/handler/mapper"
	"url-shortener-go-backend/internal/metrics"
	"url-shortener-go-backend/internal/middleware"
	"url-shortener-go-backend/internal/service"
	"url-shortener-go-backend/internal/utils"
)

type URLHandler struct {
	svc service.URLService
}

func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) HandleShorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req dto.ShortenURLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.OriginalURL) == "" {
			utils.RespondError(w, http.StatusBadRequest, "Invalid or missing URL", "")
			return
		}

		ctx := r.Context()
		userID := middleware.GetUserIDFromContext(ctx)
		var userIDPtr *string
		if userID != "" {
			userIDPtr = &userID
		}

		urlModel, err := h.svc.CreateShortURL(ctx, strings.TrimSpace(req.OriginalURL), req.IsPublic, userIDPtr, int(req.CodeLength))
		if err != nil {
			slog.Error("create short url failed", "error", err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to shorten URL", "")
			return
		}

		utils.RespondJSON(w, http.StatusCreated, mapper.ToShortenURLResponse(*urlModel), "")
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
			utils.RespondError(w, http.StatusUnauthorized, "Unauthorized", "")
			return
		}

		urls, err := h.svc.GetUserUrls(r.Context(), userID)
		if err != nil {
			slog.Error("get user urls failed", "user_id", userID, "error", err)
			utils.RespondError(w, http.StatusInternalServerError, "Could not fetch URLs", "")
			return
		}
		utils.RespondJSON(w, http.StatusOK, mapper.ToGetUrlsResponse(urls), "")
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
			utils.RespondError(w, http.StatusBadRequest, "Shortcode is required", "")
			return
		}

		ctx := r.Context()
		url, err := h.svc.GetURLByShortCode(ctx, shortcode)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.RespondError(w, http.StatusNotFound, "URL not found", "")
				return
			}
			utils.RespondError(w, http.StatusInternalServerError, "Could not fetch URL", "")
			return
		}
		if url == nil {
			utils.RespondError(w, http.StatusNotFound, "URL not found", "")
			return
		}
		utils.RespondJSON(w, http.StatusOK, mapper.ToGetURLByShortCodeResponse(*url), "")
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
			utils.RespondError(w, http.StatusBadRequest, "Missing shortcode", "")
			return
		}

		ctx := r.Context()
		urlEntry, err := h.svc.GetURLByShortCode(ctx, shortcode)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.RespondError(w, http.StatusNotFound, "URL not found", "")
				return
			}
			utils.RespondError(w, http.StatusInternalServerError, "Could not fetch URL", "")
			return
		}
		if urlEntry == nil || urlEntry.OriginalURL == "" {
			utils.RespondError(w, http.StatusNotFound, "URL not found", "")
			return
		}

		go func() {
			bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := h.svc.IncrementClickCount(bgCtx, shortcode); err != nil {
				slog.Error("click count increment failed", "shortcode", shortcode, "error", err)
			}
		}()

		metrics.URLRedirectsTotal.Inc()
		http.Redirect(w, r, urlEntry.OriginalURL, http.StatusFound)
	}
}

