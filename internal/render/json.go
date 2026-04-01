package render

import (
	"encoding/json"
	"log"
	"net/http"
)

type JSONSuccess struct {
	Success bool `json:"success"`
	Data    any  `json:"data,omitempty"`
}

type JSONError struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// JSON renders a JSON response with the given data
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	res := JSONSuccess{
		Success: true,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("ERROR: Failed to encode JSON response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ErrorJSON renders a JSON error response
func ErrorJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	res := JSONError{
		Success: false,
		Error:   message,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("ERROR: Failed to encode error response: %v", err)
	}
}
