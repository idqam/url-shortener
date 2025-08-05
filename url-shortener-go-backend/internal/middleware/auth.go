package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type contextKey string

const UserIDKey = contextKey("userID")

type JWKSCache struct {
	keySet jwk.Set
	expiry time.Time
	mutex  sync.RWMutex
}

var jwksCache = &JWKSCache{}

func AuthMiddleware(supabaseURL string) func(http.Handler) http.Handler {
	supabaseURL = strings.TrimSuffix(supabaseURL, "/")
	jwksURL := supabaseURL + "/auth/v1/.well-known/jwks.json"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			keySet, err := getJWKS(jwksURL)
			if err != nil {
				http.Error(w, "Failed to get JWKS", http.StatusInternalServerError)
				return
			}

			token, err := jwt.Parse(
				[]byte(tokenStr),
				jwt.WithKeySet(keySet),
				jwt.WithValidate(true),
			)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			sub, ok := token.Get("sub")
			if !ok {
				http.Error(w, "No subject claim", http.StatusUnauthorized)
				return
			}

			userID, ok := sub.(string)
			if !ok {
				http.Error(w, "Invalid subject claim", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getJWKS(jwksURL string) (jwk.Set, error) {
	jwksCache.mutex.RLock()
	if jwksCache.keySet != nil && time.Now().Before(jwksCache.expiry) {
		defer jwksCache.mutex.RUnlock()
		return jwksCache.keySet, nil
	}
	jwksCache.mutex.RUnlock()

	jwksCache.mutex.Lock()
	defer jwksCache.mutex.Unlock()

	if jwksCache.keySet != nil && time.Now().Before(jwksCache.expiry) {
		return jwksCache.keySet, nil
	}

	keySet, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	jwksCache.keySet = keySet
	jwksCache.expiry = time.Now().Add(time.Hour)

	return keySet, nil
}

func GetUserIDFromContext(ctx context.Context) string {
	val, _ := ctx.Value(UserIDKey).(string)
	return val
}
