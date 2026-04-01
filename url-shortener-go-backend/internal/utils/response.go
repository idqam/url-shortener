package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, data any, requestID string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if requestID != "" {
		w.Header().Set("X-Request-ID", requestID)
	}
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		if requestID != "" {
			slog.Error("json encode error", "request_id", requestID, "error", err)
		} else {
			slog.Error("json encode error", "error", err)
		}
	}
}

func RespondError(w http.ResponseWriter, status int, message string, requestID string) {
	slog.Error("response error", "status", status, "message", message)
	RespondJSON(w, status, map[string]string{"error": message}, requestID)
}
