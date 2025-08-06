package router

import (
	"context"
	"log"
	"net/http"
	"time"

	"url-shortener-go-backend/internal/handler"
	"url-shortener-go-backend/internal/middleware"
)

type APIServer struct {
	address          string
	router           *http.ServeMux
	server           *http.Server
	urlHandler       *handler.URLHandler
	analyticsHandler *handler.AnalyticsHandler
	middlewares      []func(http.Handler) http.Handler
	authMiddleware   func(http.Handler) http.Handler
}

func NewAPIServer(
	addr string,
	urlHandler *handler.URLHandler,
	analyticsHandler *handler.AnalyticsHandler,
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
	}

	s.routes()

	allMiddlewares := append(s.middlewares, s.cors(), middleware.SecurityHeaders())

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.withMiddleware(s.router, allMiddlewares...),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("[NewAPIServer] Initialized with address: %s", s.address)
	return s
}

func (s *APIServer) Run() error {
	log.Printf("[Run] HTTP server listening on %s", s.address)
	return s.server.ListenAndServe()
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	log.Println("[Shutdown] Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}

func (s *APIServer) routes() {
	log.Println("[routes] Registering HTTP routes")

	s.router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HealthCheck] %s %s", r.Method, r.URL.Path)
		handler.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	s.router.Handle("/api/urls", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.urlHandler.HandleShorten()(w, r) // Public shortening
		case http.MethodGet:
			protected := s.authMiddleware(http.HandlerFunc(s.urlHandler.HandleGetUserUrls()))
			protected.ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	s.router.HandleFunc("/api/urls/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HandleGetUrlByShortCode] %s %s", r.Method, r.URL.Path)
		s.urlHandler.HandleGetUrlByShortCode()(w, r)
	})

	s.registerAnalyticsRoutes()

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ShortCodeHandler] %s %s", r.Method, r.URL.Path)
		s.urlHandler.ShortCodeHandler()(w, r)
	})
}

func (s *APIServer) registerAnalyticsRoutes() {
	log.Println("[registerAnalyticsRoutes] Setting up analytics endpoints")

	s.router.Handle("/api/analytics/dashboard", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleGetDashboard()),
	))

	//query param
	s.router.Handle("/api/analytics/urls", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleGetTopURLs()),
	))

	s.router.Handle("/api/analytics/referrers", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleGetTopReferrers()),
	))

	s.router.Handle("/api/analytics/devices", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleGetDeviceBreakdown()),
	))

	//query param
	s.router.Handle("/api/analytics/trend", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleGetDailyTrend()),
	))

	s.router.Handle("/api/analytics/record", s.authMiddleware(
		http.HandlerFunc(s.analyticsHandler.HandleRecordAnalytics()),
	))

	log.Println("[registerAnalyticsRoutes] Analytics routes registered successfully")
}

func (s *APIServer) withMiddleware(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func (s *APIServer) cors() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			log.Printf("[CORS] %s request from origin: %s", r.Method, origin)

			w.Header().Set("Access-Control-Allow-Origin", "*") //  TODO: Restrict in production
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				log.Printf("[CORS] Preflight request handled")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
