package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"os"
)
func GenerateCode(urlStr string, length int8) (string, error) {
	if length < 6 || length > 12 {
		return "", fmt.Errorf("length must be between 6 and 12")
	}

	salt := os.Getenv("SALT")
	if salt == "" {
		return "", fmt.Errorf("missing SALT environment variable")
	}

	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	hash := sha256.Sum256(append([]byte(urlStr+salt), randomBytes...))

	numBytes := int(math.Ceil(float64(length) * 6.0 / 8.0))
	truncated := hash[:numBytes]

	encoded := base64.RawURLEncoding.EncodeToString(truncated)
	if len(encoded) < int(length) {
		return "", fmt.Errorf("unable to generate enough characters for length %d", length)
	}

	return encoded[:length], nil
}