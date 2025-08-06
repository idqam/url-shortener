package utils

import (
	"encoding/json"
	"log"
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
			log.Printf("[%s] JSON encode error: %v", requestID, err)
		} else {
			log.Printf("JSON encode error: %v", err)
		}
	}
}

func RespondError(w http.ResponseWriter, status int, message string, requestID string) {
	log.Printf("[Error %d] %s", status, message)
	RespondJSON(w, status, map[string]string{"error": message}, requestID)
}
