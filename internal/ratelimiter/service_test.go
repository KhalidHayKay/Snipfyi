package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestSaveReturnsErrorWhenCacheUnavailable(t *testing.T) {
	cache := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	svc := NewService(cache)
	l := &Limiter{conf: &Config{Every: time.Second, Rate: 1, Burst: 1}}
	if err := svc.Save(context.Background(), "key", l); err == nil {
		t.Fatal("expected error when cache is unavailable")
	}
}

func TestLoadReturnsFreshLimiterWhenCacheUnavailable(t *testing.T) {
	cache := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	svc := NewService(cache)
	conf := &Config{Every: time.Second, Rate: 1, Burst: 1}
	limiter, err := svc.Load(context.Background(), "key", conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if limiter == nil {
		t.Fatal("expected limiter")
	}
	if limiter.conf != conf {
		t.Fatal("expected returned limiter to use supplied config")
	}
}
