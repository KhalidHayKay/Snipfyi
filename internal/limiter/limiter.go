package limiter

import (
	"sync"
	"time"
)

type Limiter struct {
	Rate      float64
	Burst     float64
	Tokens    float64
	LastCheck time.Time
}

func NewLimiter(rate, burst float64) *Limiter {
	return &Limiter{
		Rate:      rate,
		Burst:     burst,
		Tokens:    burst,
		LastCheck: time.Now(),
	}
}

func (l *Limiter) Allow() bool {
	elapsed := time.Since(l.LastCheck).Seconds()
	l.Tokens += elapsed * l.Rate
	if l.Tokens > l.Burst {
		l.Tokens = l.Burst
	}
	l.LastCheck = time.Now()

	if l.Tokens >= 1 {
		l.Tokens--
		return true
	}
	return false
}

// --- middleware ---

type RateLimiter struct {
	limiters map[string]*Limiter
	mu       sync.Mutex
	Rate     float64
	Burst    float64
}

func NewRateLimiter(rate, burst float64) *RateLimiter {
	return &RateLimiter{
		Rate:  rate,
		Burst: burst,
	}
}

// func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		key := r.Header.Get("X-API-Key")
// 		if key == "" {
// 			key = utils.GetClientIP(r)
// 		}

// 		limiter, _ := loadLimiter(r.Context(), key, rl.rate, rl.burst)
// 		defer saveLimiter(r.Context(), key, limiter)

// 		if !limiter.Allow() {
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusTooManyRequests)
// 			json.NewEncoder(w).Encode(map[string]string{
// 				"error": "rate limit exceeded",
// 			})
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
