package limiter

import (
	"context"
	"smply/internal/storage"
	"strconv"
	"time"
)

func saveLimiter(ctx context.Context, key string, l *Limiter) error {
	return storage.Cache.HSet(ctx, "ratelimit:"+key, map[string]interface{}{
		"rate":      l.rate,
		"burst":     l.burst,
		"tokens":    l.tokens,
		"lastCheck": l.lastCheck.UnixMilli(),
	}).Err()
}

func loadLimiter(ctx context.Context, key string, rate, burst float64) (*Limiter, error) {
	vals, err := storage.Cache.HGetAll(ctx, "ratelimit:"+key).Result()
	if err != nil || len(vals) == 0 {
		// Key doesn't exist yet, create fresh
		return NewLimiter(rate, burst), nil
	}

	tokens, _ := strconv.ParseFloat(vals["tokens"], 64)
	lastCheckMs, _ := strconv.ParseInt(vals["lastCheck"], 10, 64)

	return &Limiter{
		rate:      rate, // or parse from hash if you want per-key rates
		burst:     burst,
		tokens:    tokens,
		lastCheck: time.UnixMilli(lastCheckMs),
	}, nil
}
