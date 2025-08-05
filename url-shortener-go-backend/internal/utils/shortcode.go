package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func GenerateCode(urlStr string, length int) (string, error) {
	if length < 6 || length > 12 {
		return "", fmt.Errorf("length must be between 6 and 12")
	}

	if strings.TrimSpace(urlStr) == "" {
		return "", fmt.Errorf("urlStr cannot be empty")
	}

	salt := os.Getenv("SALT")
	if salt == "" {
		return "", fmt.Errorf("missing SALT environment variable")
	}

	if len(salt) < 32 {
		return "", fmt.Errorf("SALT must be at least 32 characters long")
	}

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	timestamp := time.Now().UnixNano()
	timestampBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		timestampBytes[i] = byte(timestamp >> (i * 8))
	}

	hasher := sha256.New()
	hasher.Write([]byte(urlStr))
	hasher.Write([]byte(salt))
	hasher.Write(randomBytes)
	hasher.Write(timestampBytes)
	hash := hasher.Sum(nil)

	numBytes := int(math.Ceil(float64(length) * 6.0 / 8.0))
	if numBytes > len(hash) {
		numBytes = len(hash)
	}

	encoded := base64.RawURLEncoding.EncodeToString(hash[:numBytes])

	if len(encoded) < length {
		additionalBytes := make([]byte, 8)
		rand.Read(additionalBytes)
		additionalHash := sha256.Sum256(append(hash, additionalBytes...))
		additionalEncoded := base64.RawURLEncoding.EncodeToString(additionalHash[:])
		encoded += additionalEncoded
	}

	result := encoded[:length]

	return result, nil
}
