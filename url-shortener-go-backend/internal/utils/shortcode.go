package utils

import (
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

    hash := sha256.Sum256([]byte(urlStr + salt))
    numBytes := int(math.Ceil(float64(length) * 6.0 / 8.0))
    truncated := hash[:numBytes]

    encoded := base64.RawURLEncoding.EncodeToString(truncated)
    if len(encoded) < int(length) {
        return "", fmt.Errorf("unable to generate enough characters for length %d", length)
    }

    return encoded[:length], nil
}
