package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"smply/app/render"
	"smply/config"
	"smply/internal/apikey"
	"smply/internal/limiter"
	"smply/internal/session"
	"smply/utils"
	"strings"
)

type Middleware struct {
	apikeyService  *apikey.Service
	sessionService *session.Service
	limiterService *limiter.Service
}

func NewMiddleware(
	apikeyService *apikey.Service,
	sessionService *session.Service,
	limiterService *limiter.Service,
) *Middleware {
	return &Middleware{
		apikeyService:  apikeyService,
		sessionService: sessionService,
		limiterService: limiterService,
	}
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

		if err != nil || !m.sessionService.ValidateSession(r.Context(), sessionId.Value) {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) AdminGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, err := r.Cookie(strings.ToLower(config.Env.App.Name) + "_session_id")

		if err == nil && m.sessionService.ValidateSession(r.Context(), sessionId.Value) {
			http.Redirect(w, r, "/admin", http.StatusFound)
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RateLimit(rl *limiter.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			key = utils.GetClientIP(r)
		}

		limiter, _ := m.limiterService.Load(r.Context(), key, rl.Rate, rl.Burst)
		defer m.limiterService.Save(r.Context(), key, limiter)

		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "rate limit exceeded",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
