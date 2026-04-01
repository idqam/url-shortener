package router

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"url-shortener-go-backend/internal/cache"
	"url-shortener-go-backend/internal/config"
	"url-shortener-go-backend/internal/handler"
	"url-shortener-go-backend/internal/metrics"
	"url-shortener-go-backend/internal/middleware"
	"url-shortener-go-backend/internal/repository"
	"url-shortener-go-backend/internal/utils"
)

type APIServer struct {
	address          string
	router           *http.ServeMux
	server           *http.Server
	urlHandler       *handler.URLHandler
	analyticsHandler *handler.AnalyticsHandler
	middlewares      []func(http.Handler) http.Handler
	authMiddleware   func(http.Handler) http.Handler
	cache            cache.Cache
	supabaseRepo     *repository.SupabaseRepository
	cfg              *config.Config
}

func NewAPIServer(
	addr string,
	cfg *config.Config,
	urlHandler *handler.URLHandler,
	analyticsHandler *handler.AnalyticsHandler,
	c cache.Cache,
	supabaseRepo *repository.SupabaseRepository,
	authMw func(http.Handler) http.Handler,
	mws ...func(http.Handler) http.Handler,
) *APIServer {
	mux := http.NewServeMux()
	s := &APIServer{
		address:          addr,
		router:           mux,
		urlHandler:       urlHandler,
		analyticsHandler: analyticsHandler,
		middlewares:      mws,
		authMiddleware:   authMw,
		cache:            c,
		supabaseRepo:     supabaseRepo,
		cfg:              cfg,
	}

	s.routes()

	isDev := cfg.Environment == "development"
	allMiddlewares := append(s.middlewares,
		middleware.RequestIDMiddleware,
		middleware.MetricsMiddleware,
		middleware.TracingMiddleware,
		s.cors(),
		middleware.SecurityHeaders(isDev),
	)

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.withMiddleware(s.router, allMiddlewares...),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("api server initialized", "address", s.address)
	return s
}

func (s *APIServer) Run() error {
	slog.Info("http server listening", "address", s.address)
	return s.server.ListenAndServe()
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	slog.Info("shutting down http server")
	return s.server.Shutdown(ctx)
}

func (s *APIServer) routes() {
	slog.Info("registering http routes")

	s.router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("health check", "method", r.Method, "path", r.URL.Path)
		ctx := r.Context()

		redisStatus := "connected"
		if s.cache != nil {
			if err := s.cache.Ping(ctx); err != nil {
				redisStatus = "disconnected"
				slog.Warn("redis ping failed during health check", "error", err)
			}
		} else {
			redisStatus = "disabled"
		}

		dbStatus := "connected"
		if s.supabaseRepo != nil {
			_, _, err := s.supabaseRepo.Client.From("urls").Select("id", "exact", false).Limit(0, "").Execute()
			if err != nil {
				dbStatus = "disconnected"
				slog.Warn("supabase ping failed during health check", "error", err)
			}
		}

		overallStatus := "ok"
		if redisStatus == "disconnected" || dbStatus == "disconnected" {
			overallStatus = "degraded"
		}

		utils.RespondJSON(w, http.StatusOK, map[string]string{
			"status":   overallStatus,
			"redis":    redisStatus,
			"database": dbStatus,
			"version":  s.cfg.Version,
		}, "")
	})

	s.router.Handle("/metrics", metrics.Handler())

	s.router.Handle("/api/urls", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.urlHandler.HandleShorten()(w, r)
		case http.MethodGet:
			protected := s.authMiddleware(http.HandlerFunc(s.urlHandler.HandleGetUserUrls()))
			protected.ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	s.router.HandleFunc("/api/urls/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("get url by shortcode", "method", r.Method, "path", r.URL.Path)
		s.urlHandler.HandleGetUrlByShortCode()(w, r)
	})

	s.registerAnalyticsRoutes()

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("shortcode handler", "method", r.Method, "path", r.URL.Path)
		s.urlHandler.ShortCodeHandler()(w, r)
	})
}

func (s *APIServer) registerAnalyticsRoutes() {
	slog.Info("registering analytics routes")

	s.router.Handle("/api/analytics/dashboard", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleGetDashboard()),
	))

	s.router.Handle("/api/analytics/urls", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleGetTopURLs()),
	))

	s.router.Handle("/api/analytics/referrers", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleGetTopReferrers()),
	))

	s.router.Handle("/api/analytics/devices", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleGetDeviceBreakdown()),
	))

	s.router.Handle("/api/analytics/trend", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleGetDailyTrend()),
	))

	s.router.Handle("/api/analytics/record", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleRecordAnalytics()),
	))

	slog.Info("analytics routes registered")
}

func (s *APIServer) withMiddleware(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func (s *APIServer) cors() func(http.Handler) http.Handler {
	allowedOrigins := s.cfg.AllowedOrigins
	allowedSet := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowedSet[strings.TrimRight(o, "/")] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			slog.Debug("cors check", "method", r.Method, "origin", origin)

			if _, ok := allowedSet[strings.TrimRight(origin, "/")]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				slog.Debug("cors preflight handled")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
