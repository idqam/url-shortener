package middleware

import (
	"context"
	"net"
	"net/http"
	"time"

	"url-shortener-go-backend/internal/cache"
)

type RateLimiter struct {
	cache     cache.Cache
	limit     int
	windowTTL time.Duration
}

func NewRateLimiter(c cache.Cache, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{cache: c, limit: limit, windowTTL: window}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var key string
		userID := GetUserIDFromContext(ctx)
		if userID != "" {
			key = "ratelimit:user:" + userID
		} else {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			key = "ratelimit:ip:" + ip
		}

		count, err := rl.increment(ctx, key)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if count > rl.limit {
			http.Error(w, "Rate limit exceeded — slow down", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) increment(ctx context.Context, key string) (int, error) {
	val, err := rl.cache.Incr(ctx, key)
	if err != nil {
		return 0, err
	}
	if val == 1 {
		_ = rl.cache.Expire(ctx, key, rl.windowTTL)
	}
	return int(val), nil
}
