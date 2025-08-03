package utils

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func ExtractEmailFromSupabaseToken(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse JWT: %w", err)
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

func ExtractUserIDFromSupabaseToken(tokenString string) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok && len(sub) > 0 {
    return sub, nil
}

		return "", errors.New("missing or invalid `sub` claim")
	}

	return "", errors.New("invalid token claims")
}
