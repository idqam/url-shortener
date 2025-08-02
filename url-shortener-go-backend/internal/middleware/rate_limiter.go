package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"url-shortener-go-backend/internal/cache"

	"github.com/golang-jwt/jwt"
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

		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := rl.extractUserIDFromToken(token)
			if err == nil && userID != "" {
				key = "ratelimit:user:" + userID
			}
		}

		if key == "" {
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

func (rl *RateLimiter) extractUserIDFromToken(token string) (string, error) {
	claims, err := ExtractClaimsFromToken(token)
	if err != nil {

		return "", nil
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", nil
	}
	return sub, nil
}

func ExtractClaimsFromToken(tokenStr string) (map[string]interface{}, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid claims type")
}
