package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	secret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secret: secret}
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		email, err := am.extractEmailFromToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_email", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) extractEmailFromToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(am.secret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("token parse error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", errors.New("email not found in token")
	}

	return email, nil
}
