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
			log.Println("[WARN] Invalid REDIS_URL format. Cache disabled.")
		} else {
			rc = cache.NewRedisCache(addr, password, 0, 24*time.Hour)
			log.Printf("[INFO] Redis cache enabled")
		}
	} else {
		log.Println("[INFO] REDIS_URL not provided. Cache is disabled.")
	}

	apiKey := os.Getenv("SERVICE_ROLE")
	apiURL := os.Getenv("DB_API_URL")
	if apiKey == "" || apiURL == "" {
		log.Fatal("[FATAL] Required environment variables are missing")
	}

	supabase, err := repository.NewSupabaseRepository(apiURL, apiKey)
	if err != nil {
		log.Fatal("[FATAL] Failed to initialize database connection")
	}

	jwtSecret := os.Getenv("DB_API_URL")
	if jwtSecret == "" {
		log.Fatal("[FATAL] JWT_SECRET environment variable is required")
	}

	authMw := middleware.AuthMiddleware(jwtSecret)
	urlRepo := repository.NewURLRepository(supabase)
	analyticsRepo := repository.NewAnalyticsRepository(supabase)
	urlService := service.NewURLService(urlRepo, rc)
	urlHandler := handler.NewURLHandler(urlService, rc)
	analyticsService := service.NewAnalyticsService(analyticsRepo, rc)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	var rateLimiterConfig *middleware.RateLimiterConfig
	if os.Getenv("ENVIRONMENT") == "development" {
		rateLimiterConfig = middleware.DevelopmentRateLimiterConfig()
	} else {
		rateLimiterConfig = middleware.ProductionRateLimiterConfig()
	}
	limiter := middleware.NewRateLimiter(rc, rateLimiterConfig)

	port := os.Getenv("PORT")
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
			log.Fatal("[FATAL] Server failed to start")
		}
	}()

	log.Printf("[INFO] Server listening on port %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("[ERROR] Server shutdown timed out")
		os.Exit(1)
	}

	log.Println("[INFO] Server stopped")
}

func parseRedisURL(raw string) (addr, password string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}
	pw, _ := u.User.Password()
	return u.Host, pw, nil
}
