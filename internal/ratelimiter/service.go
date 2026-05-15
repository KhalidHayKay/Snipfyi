package ratelimiter

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	cache *redis.Client
}

func NewService(cache *redis.Client) *Service {
	return &Service{
		cache: cache,
	}
}

func (s *Service) Save(ctx context.Context, key string, l *Limiter) error {
	return s.cache.HSet(ctx, "ratelimit:"+key, map[string]interface{}{
		"tokens":    l.tokens,
		"lastCheck": l.lastCheck.UnixMilli(),
	}).Err()
}

func (s *Service) Load(ctx context.Context, key string, conf *Config) (*Limiter, error) {
	vals, err := s.cache.HGetAll(ctx, "ratelimit:"+key).Result()
	if err != nil || len(vals) == 0 {
		// Key doesn't exist yet, create fresh
		return NewLimiter(conf), nil
	}

	tokens, _ := strconv.ParseFloat(vals["tokens"], 64)
	lastCheckMs, _ := strconv.ParseInt(vals["lastCheck"], 10, 64)

	return &Limiter{
		conf:      conf,
		tokens:    tokens,
		lastCheck: time.UnixMilli(lastCheckMs),
	}, nil
}
