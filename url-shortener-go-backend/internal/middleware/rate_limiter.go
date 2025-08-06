package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"url-shortener-go-backend/internal/cache"
)

type RateLimiter struct {
	cache     cache.Cache
	config    *RateLimiterConfig
	whitelist map[string]bool
	mu        sync.RWMutex
}

type RateLimiterConfig struct {
	DefaultLimit int
	WindowTTL    time.Duration

	// Tiered limits
	AnonymousLimit     int
	AuthenticatedLimit int
	PremiumLimit       int

	// Burst configuration
	BurstEnabled    bool
	BurstMultiplier float64
	BurstDuration   time.Duration

	// Response configuration
	IncludeHeaders bool
	CustomMessage  string

	DistributedMode bool
	SlidingWindow   bool
}

func NewRateLimiter(c cache.Cache, config *RateLimiterConfig) *RateLimiter {
	if config == nil {
		config = DefaultRateLimiterConfig()
	}

	return &RateLimiter{
		cache:     c,
		config:    config,
		whitelist: make(map[string]bool),
	}
}

func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		DefaultLimit:       100,
		WindowTTL:          time.Minute,
		AnonymousLimit:     50,
		AuthenticatedLimit: 100,
		PremiumLimit:       500,
		BurstEnabled:       true,
		BurstMultiplier:    1.5,
		BurstDuration:      10 * time.Second,
		IncludeHeaders:     true,
		CustomMessage:      "Rate limit exceeded. Please slow down.",
		DistributedMode:    false,
		SlidingWindow:      false,
	}
}
func ProductionRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		DefaultLimit: 100,
		WindowTTL:    time.Minute,

		AnonymousLimit:     20,  // Limited for non-authenticated
		AuthenticatedLimit: 100, // Standard for logged-in users
		PremiumLimit:       500, // Higher for premium users

		// Burst configuration (handle traffic spikes)
		BurstEnabled:    true,
		BurstMultiplier: 1.5,
		BurstDuration:   30 * time.Second,

		IncludeHeaders: true,
		CustomMessage:  "Rate limit exceeded. Please try again later.",

		DistributedMode: true, // Important for production
		SlidingWindow:   false,
	}
}

func DevelopmentRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		DefaultLimit:       1000,
		WindowTTL:          time.Minute,
		AnonymousLimit:     1000,
		AuthenticatedLimit: 1000,
		PremiumLimit:       1000,

		BurstEnabled:   false,
		IncludeHeaders: true,
		CustomMessage:  "Rate limit exceeded (dev mode)",

		DistributedMode: false,
		SlidingWindow:   false,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		identifier := rl.getIdentifier(ctx, r)

		if rl.isWhitelisted(identifier.IP) {
			next.ServeHTTP(w, r)
			return
		}

		limit := rl.getLimit(identifier)

		key := rl.buildKey(identifier)

		result, err := rl.checkRateLimit(ctx, key, limit)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if rl.config.IncludeHeaders {
			rl.setHeaders(w, result)
		}

		if result.Exceeded {
			rl.handleLimitExceeded(w, r, result)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) CustomMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			identifier := rl.getIdentifier(ctx, r)

			key := fmt.Sprintf("%s:custom:%s", rl.buildKey(identifier), r.URL.Path)

			result, err := rl.checkRateLimitWithCustom(ctx, key, limit, window)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if rl.config.IncludeHeaders {
				rl.setHeaders(w, result)
			}

			if result.Exceeded {
				rl.handleLimitExceeded(w, r, result)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type Identifier struct {
	UserID string
	IP     string
	Tier   string // "anonymous", "authenticated", "premium"
}

type RateLimitResult struct {
	Count     int
	Limit     int
	Remaining int
	ResetAt   time.Time
	Exceeded  bool
}

func (rl *RateLimiter) getIdentifier(ctx context.Context, r *http.Request) Identifier {
	userID := GetUserIDFromContext(ctx)
	ip := rl.extractIP(r)

	tier := "anonymous"
	if userID != "" {
		tier = "authenticated"

	}

	return Identifier{
		UserID: userID,
		IP:     ip,
		Tier:   tier,
	}
}

func (rl *RateLimiter) extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return strings.TrimSpace(cfIP)
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func (rl *RateLimiter) getLimit(identifier Identifier) int {
	switch identifier.Tier {
	case "premium":
		return rl.config.PremiumLimit
	case "authenticated":
		return rl.config.AuthenticatedLimit
	case "anonymous":
		return rl.config.AnonymousLimit
	default:
		return rl.config.DefaultLimit
	}
}

func (rl *RateLimiter) buildKey(identifier Identifier) string {
	if identifier.UserID != "" {
		return fmt.Sprintf("ratelimit:user:%s", identifier.UserID)
	}
	return fmt.Sprintf("ratelimit:ip:%s", identifier.IP)
}

func (rl *RateLimiter) checkRateLimit(ctx context.Context, key string, limit int) (*RateLimitResult, error) {
	if rl.config.SlidingWindow {
		return rl.checkSlidingWindow(ctx, key, limit)
	}
	return rl.checkFixedWindow(ctx, key, limit)
}

func (rl *RateLimiter) checkFixedWindow(ctx context.Context, key string, limit int) (*RateLimitResult, error) {
	if rl.config.BurstEnabled {
		burstKey := key + ":burst"
		burstLimit := int(float64(limit) * rl.config.BurstMultiplier)

		burstCount, _, err := rl.cache.Get(ctx, burstKey)
		if err == nil && burstCount != "" {
			limit = burstLimit
		}
	}

	count, err := rl.cache.Incr(ctx, key)
	if err != nil {
		return nil, err
	}

	if count == 1 {
		if err := rl.cache.Expire(ctx, key, rl.config.WindowTTL); err != nil {
			return nil, err
		}
	}

	ttl := rl.config.WindowTTL

	result := &RateLimitResult{
		Count:     int(count),
		Limit:     limit,
		Remaining: max(0, limit-int(count)),
		ResetAt:   time.Now().Add(ttl),
		Exceeded:  int(count) > limit,
	}

	if result.Exceeded && rl.config.BurstEnabled {
		burstKey := key + ":burst"
		_ = rl.cache.Set(ctx, burstKey, "1", rl.config.BurstDuration)
	}

	return result, nil
}

func (rl *RateLimiter) checkSlidingWindow(ctx context.Context, key string, limit int) (*RateLimitResult, error) {

	return rl.checkFixedWindow(ctx, key, limit)
}

func (rl *RateLimiter) checkRateLimitWithCustom(ctx context.Context, key string, limit int, window time.Duration) (*RateLimitResult, error) {
	count, err := rl.cache.Incr(ctx, key)
	if err != nil {
		return nil, err
	}

	if count == 1 {
		if err := rl.cache.Expire(ctx, key, window); err != nil {
			return nil, err
		}
	}

	return &RateLimitResult{
		Count:     int(count),
		Limit:     limit,
		Remaining: max(0, limit-int(count)),
		ResetAt:   time.Now().Add(window),
		Exceeded:  int(count) > limit,
	}, nil
}

func (rl *RateLimiter) setHeaders(w http.ResponseWriter, result *RateLimitResult) {
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", result.Limit))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", result.ResetAt.Unix()))

	if result.Exceeded {
		w.Header().Set("Retry-After", fmt.Sprintf("%d", int(time.Until(result.ResetAt).Seconds())))
	}
}

func (rl *RateLimiter) handleLimitExceeded(w http.ResponseWriter, r *http.Request, result *RateLimitResult) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)

	response := map[string]interface{}{
		"error":       rl.config.CustomMessage,
		"retry_after": int(time.Until(result.ResetAt).Seconds()),
		"limit":       result.Limit,
	}

	json.NewEncoder(w).Encode(response)
}

func (rl *RateLimiter) AddToWhitelist(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.whitelist[ip] = true
}

func (rl *RateLimiter) RemoveFromWhitelist(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.whitelist, ip)
}

func (rl *RateLimiter) isWhitelisted(ip string) bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.whitelist[ip]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
