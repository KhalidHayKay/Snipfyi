package middleware

import (
	"log"
	"net/http"
	"smply/config"
	"smply/internal/render"
	"smply/internal/service"
	"strings"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

func AdminAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, err := r.Cookie(strings.ToLower(config.Env.App.Name) + "_session_id")

		if err != nil || !service.ValidateSession(r.Context(), sessionId.Value) {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
		}

		next.ServeHTTP(w, r)
	})
}

func AdminGuestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, err := r.Cookie(strings.ToLower(config.Env.App.Name) + "_session_id")

		if err == nil && service.ValidateSession(r.Context(), sessionId.Value) {
			http.Redirect(w, r, "/admin", http.StatusFound)
		}

		next.ServeHTTP(w, r)
	})
}
