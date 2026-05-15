package session

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestCreateSessionReturnsErrorWhenCacheUnavailable(t *testing.T) {
	cache := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	svc := NewService(cache)
	_, err := svc.CreateSession(context.Background(), "identity")
	if err == nil {
		t.Fatal("expected error when cache is unavailable")
	}
}

func TestValidateSessionReturnsFalseWhenCacheUnavailable(t *testing.T) {
	cache := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	svc := NewService(cache)
	if svc.ValidateSession(context.Background(), "session123") {
		t.Fatal("expected false when cache is unavailable")
	}
}

func TestDeleteSessionReturnsFalseWhenCacheUnavailable(t *testing.T) {
	cache := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	svc := NewService(cache)
	if svc.DeleteSession(context.Background(), "session123") {
		t.Fatal("expected false when cache is unavailable")
	}
}
