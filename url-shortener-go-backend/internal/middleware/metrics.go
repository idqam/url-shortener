package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"url-shortener-go-backend/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.HTTPRequestsInFlight.Inc()
		defer metrics.HTTPRequestsInFlight.Dec()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()

		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()
		path := normalizePath(r.URL.Path)

		metrics.HTTPRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, path, fmt.Sprintf("%d", rec.status)).Inc()
	})
}

func normalizePath(path string) string {
	if strings.HasPrefix(path, "/api/") {
		return path
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 1 && len(parts[0]) > 0 {
		return "/:shortcode"
	}
	return path
}
