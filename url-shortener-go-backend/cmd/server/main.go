package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener-go-backend/internal/cache"
	"url-shortener-go-backend/internal/config"
	"url-shortener-go-backend/internal/handler"
	"url-shortener-go-backend/internal/logger"
	"url-shortener-go-backend/internal/metrics"
	"url-shortener-go-backend/internal/middleware"
	"url-shortener-go-backend/internal/repository"
	"url-shortener-go-backend/internal/router"
	"url-shortener-go-backend/internal/service"
	"url-shortener-go-backend/internal/telemetry"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Info("no .env file found, using system environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Init(cfg.Environment)

	slog.Info("starting server", "port", cfg.Port, "env", cfg.Environment)

	shutdownTracer, err := telemetry.InitTracer("url-shortener", cfg.Environment, cfg.OTLPEndpoint)
	if err != nil {
		slog.Error("failed to initialize tracer", "error", err)
		os.Exit(1)
	}

	metrics.Init()

	var rc cache.Cache
	if cfg.RedisURL != "" {
		addr, password, err := parseRedisURL(cfg.RedisURL)
		if err != nil {
			slog.Warn("invalid REDIS_URL format, cache disabled", "error", err)
		} else {
			rc = cache.NewRedisCache(addr, password, 0, 24*time.Hour)
			slog.Info("redis cache enabled")
		}
	} else {
		slog.Info("REDIS_URL not provided, cache disabled")
	}

	if rc != nil {
		rc = cache.NewInstrumentedCache(rc)
	}

	supabase, err := repository.NewSupabaseRepository(cfg.DBAPIUrl, cfg.SupabaseServiceRole)
	if err != nil {
		slog.Error("failed to initialize database connection", "error", err)
		os.Exit(1)
	}

	urlRepo := repository.NewURLRepository(supabase, cfg.ShortDomain)
	analyticsRepo := repository.NewAnalyticsRepository(supabase)

	urlRepo = repository.NewInstrumentedURLRepository(urlRepo)
	analyticsRepo = repository.NewInstrumentedAnalyticsRepository(analyticsRepo)

	urlService := service.NewURLService(urlRepo, rc, cfg.Salt)
	analyticsService := service.NewAnalyticsService(analyticsRepo, rc, cfg.Salt)

	urlHandler := handler.NewURLHandler(urlService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	authMw := middleware.AuthMiddleware(cfg.JWTSecret)

	var rateLimiterConfig *middleware.RateLimiterConfig
	if cfg.Environment == "development" {
		rateLimiterConfig = middleware.DevelopmentRateLimiterConfig()
	} else {
		rateLimiterConfig = middleware.ProductionRateLimiterConfig()
	}
	limiter := middleware.NewRateLimiter(rc, rateLimiterConfig)

	server := router.NewAPIServer(
		":"+cfg.Port,
		cfg,
		urlHandler,
		analyticsHandler,
		rc,
		supabase,
		authMw,
		limiter.Middleware,
	)

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("server listening", "port", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown timed out", "error", err)
		os.Exit(1)
	}

	if rc != nil {
		if err := rc.Close(); err != nil {
			slog.Error("failed to close redis connection", "error", err)
		}
	}

	tracerCtx, tracerCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer tracerCancel()
	if err := shutdownTracer(tracerCtx); err != nil {
		slog.Error("failed to shutdown tracer", "error", err)
	}

	slog.Info("server stopped")
}

func parseRedisURL(raw string) (addr, password string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}
	pw, _ := u.User.Password()
	return u.Host, pw, nil
}
