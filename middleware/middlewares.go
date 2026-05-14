package middleware

import (
	"log"
	"net/http"
	"smply/config"
	"smply/internal/apikey"
	"smply/internal/render"
	"smply/internal/service"
	"strings"
)

type Middleware struct {
	apikeyService *apikey.Service
}

func NewMiddleware(apikeyService *apikey.Service) *Middleware {
	return &Middleware{apikeyService}
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
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

func (m *Middleware) RequireKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			render.ErrorJSON(w, http.StatusUnauthorized, "API key required")
			return
		}

		valid, err := m.apikeyService.Validate(r.Context(), key)
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

func (m *Middleware) AuthenticateAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, err := r.Cookie(strings.ToLower(config.Env.App.Name) + "_session_id")

		if err != nil || !service.ValidateSession(r.Context(), sessionId.Value) {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) AdminGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, err := r.Cookie(strings.ToLower(config.Env.App.Name) + "_session_id")

		if err == nil && service.ValidateSession(r.Context(), sessionId.Value) {
			http.Redirect(w, r, "/admin", http.StatusFound)
		}

		next.ServeHTTP(w, r)
	})
}
