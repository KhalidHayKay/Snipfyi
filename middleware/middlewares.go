package middleware

import (
	"log"
	"net/http"
	"smply/internal/render"
	"smply/service"
)

func RequireKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			render.ErrorJSON(w, http.StatusUnauthorized, "API key required")
			return
		}

		valid, err := service.ValidateAPIKey(r.Context(), key)
		if err != nil {
			log.Println("Error validating API key:", err)
			render.ErrorJSON(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		if !valid {
			render.ErrorJSON(w, http.StatusUnauthorized, "Invalid or expired API key")
			return
		}

		next.ServeHTTP(w, r)
	})
}
