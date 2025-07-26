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
}

func NewAPIServer(addr string, h *handler.URLHandler) *APIServer {
	mux := http.NewServeMux()
	s := &APIServer{
		address: addr,
		router:  mux,
		h:       h,
	}

	s.routes()
	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.withMiddleware(s.router, s.cors()),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *APIServer) Run() error {
	log.Printf("HTTP server listening on %s", s.address)
	return s.server.ListenAndServe()
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}

func (s *APIServer) routes() {
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	s.router.HandleFunc("/urls", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.h.HandleShorten()(w, r)
		case http.MethodGet:
			s.h.HandleGetUserUrls()(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	s.router.HandleFunc("/url", s.h.HandleGetUrlByShortCode())
	s.router.HandleFunc("/", s.h.HandleRedirect())
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
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
