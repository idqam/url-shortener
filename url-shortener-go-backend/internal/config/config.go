package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Port                string
	Environment         string
	AllowedOrigins      []string
	RedisURL            string
	DBAPIUrl            string
	SupabaseServiceRole string
	JWTSecret           string
	Salt                string
	ShortDomain         string
	ShutdownTimeout     time.Duration
	OTLPEndpoint        string
	Version             string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "production"
	}

	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		for _, o := range strings.Split(allowedOriginsStr, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				allowedOrigins = append(allowedOrigins, o)
			}
		}
	}
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"https://url-shortener-nu-two-32.vercel.app"}
	}

	dbAPIUrl := os.Getenv("DB_API_URL")
	if dbAPIUrl == "" {
		return nil, fmt.Errorf("DB_API_URL is required")
	}

	serviceRole := os.Getenv("SERVICE_ROLE")
	if serviceRole == "" {
		return nil, fmt.Errorf("SERVICE_ROLE is required")
	}

	jwtSecret := os.Getenv("SUPABASE_URL")
	if jwtSecret == "" {
		return nil, fmt.Errorf("SUPABASE_URL is required")
	}

	salt := os.Getenv("SALT")
	if salt == "" {
		return nil, fmt.Errorf("SALT is required")
	}

	shortDomain := os.Getenv("SHORT_DOMAIN")
	if shortDomain == "" {
		shortDomain = "http://localhost:8080"
	}

	shutdownTimeout := 15 * time.Second
	if env == "development" {
		shutdownTimeout = 5 * time.Second
	}

	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "dev"
	}

	return &Config{
		Port:                port,
		Environment:         env,
		AllowedOrigins:      allowedOrigins,
		RedisURL:            os.Getenv("REDIS_URL"),
		DBAPIUrl:            dbAPIUrl,
		SupabaseServiceRole: serviceRole,
		JWTSecret:           jwtSecret,
		Salt:                salt,
		ShortDomain:         shortDomain,
		ShutdownTimeout:     shutdownTimeout,
		OTLPEndpoint:        os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		Version:             version,
	}, nil
}
