package router

import (
	"context"
	"log"
	"net/http"
	"time"

	"url-shortener-go-backend/internal/handler"
)

type APIServer struct {
	address string
	router  *http.ServeMux
	server  *http.Server
	h       *handler.URLHandler
	us *handler.UserHandler
}

func NewAPIServer(addr string, h *handler.URLHandler, us *handler.UserHandler) *APIServer {
	mux := http.NewServeMux()
	s := &APIServer{
		address: addr,
		router:  mux,
		h:       h,
		us: us,
	}

	s.routes()

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.withMiddleware(s.router, s.cors()),
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
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	s.router.HandleFunc("/api/urls", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[/api/urls] %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodPost:
			s.h.HandleShorten()(w, r)
		case http.MethodGet:
			s.h.HandleGetUserUrls()(w, r)
		default:
			log.Printf("[/api/urls] Method not allowed: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	s.router.HandleFunc("/api/url", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[/api/url] %s %s", r.Method, r.URL.Path)
		s.h.HandleGetUrlByShortCode()(w, r)
	})

	s.router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
	log.Printf("[/api/users] %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	s.us.UserCreationHandler()(w, r)
})

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ShortCodeHandler] %s %s", r.Method, r.URL.Path)
		s.h.ShortCodeHandler()(w, r)
	})
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

			w.Header().Set("Access-Control-Allow-Origin", "*") // 🔒 TODO: restrict in production
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
