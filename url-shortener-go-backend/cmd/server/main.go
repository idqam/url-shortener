package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener-go-backend/internal/cache"
	"url-shortener-go-backend/internal/handler"
	"url-shortener-go-backend/internal/middleware"
	"url-shortener-go-backend/internal/repository"
	"url-shortener-go-backend/internal/router"
	"url-shortener-go-backend/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("[INFO] No .env file found. Using system environment variables.")
	}

	var rc cache.Cache
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		addr, password, err := parseRedisURL(redisURL)
		if err != nil {
			log.Printf("[WARN] Invalid REDIS_URL. Cache disabled: %v", err)
		} else {
			rc = cache.NewRedisCache(addr, password, 0, 24*time.Hour)
			log.Printf("[INFO] Redis cache enabled at %s", addr)
		}
	} else {
		log.Println("[INFO] REDIS_URL not provided. Cache is disabled.")
	}

	apiKey := os.Getenv("SERVICE_ROLE")
	apiURL := os.Getenv("DB_API_URL")
	if apiKey == "" || apiURL == "" {
		log.Fatal("[FATAL] SERVICE_ROLE or DB_API_URL environment variable is missing")
	}

	supabase, err := repository.NewSupabaseRepository(apiURL, apiKey)
	if err != nil {
		log.Fatalf("[FATAL] Failed to create Supabase repository: %v", err)
	}

	jwtSecret := os.Getenv("DB_API_URL")
	if jwtSecret == "" {
		log.Fatal("[FATAL] DB_API_URL is not set")
	}

	authMw := middleware.AuthMiddleware(jwtSecret)

	urlRepo := repository.NewURLRepository(supabase)
	analyticsRepo := repository.NewAnalyticsRepository(supabase)
	urlService := service.NewURLService(urlRepo, rc)
	urlHandler := handler.NewURLHandler(urlService, rc)
	analyticsService := service.NewAnalyticsService(analyticsRepo, rc)

	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	limiter := middleware.NewRateLimiter(rc, 10, 1*time.Minute)
	port := os.Getenv("PORT") //change in prod
	if port == "" {
		port = "8080"
	}

	server := router.NewAPIServer(
		":"+port,
		urlHandler,
		analyticsHandler,
		authMw,
		limiter.Middleware,
	)

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] HTTP server error: %v", err)
		}
	}()
	log.Printf("[INFO] HTTP server listening on :%s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("[INFO] Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Server shutdown failed: %v", err)
	}
	log.Println("[INFO] Server exited cleanly")
}

func parseRedisURL(raw string) (addr, password string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}
	pw, _ := u.User.Password()
	return u.Host, pw, nil
}
