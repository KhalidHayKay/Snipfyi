package limiter

import (
	"encoding/json"
	"net/http"
	"smply/utils"
	"sync"
	"time"
)

type Limiter struct {
	rate      float64
	burst     float64
	tokens    float64
	lastCheck time.Time
}

func NewLimiter(rate, burst float64) *Limiter {
	return &Limiter{
		rate:      rate,
		burst:     burst,
		tokens:    burst,
		lastCheck: time.Now(),
	}
}

func (l *Limiter) Allow() bool {
	elapsed := time.Since(l.lastCheck).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > l.burst {
		l.tokens = l.burst
	}
	l.lastCheck = time.Now()

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

// --- middleware ---

type RateLimiter struct {
	limiters map[string]*Limiter
	mu       sync.Mutex
	rate     float64
	burst    float64
}

func NewRateLimiter(rate, burst float64) *RateLimiter {
	return &RateLimiter{
		rate:  rate,
		burst: burst,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			key = utils.GetClientIP(r)
		}

		limiter, _ := loadLimiter(r.Context(), key, rl.rate, rl.burst)
		defer saveLimiter(r.Context(), key, limiter)

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
