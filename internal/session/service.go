package session

import (
	"context"
	"smply/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	cache *redis.Client
}

func NewService(cache *redis.Client) *Service {
	return &Service{cache}
}

func (s *Service) CreateSession(ctx context.Context, identity string) (string, error) {
	sessionId, err := utils.GenerateToken()
	if err != nil {
		return "", err
	}

	err = s.cache.Set(ctx, sessionId, identity, 1*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *Service) ValidateSession(ctx context.Context, sessionId string) bool {
	err := s.cache.Get(ctx, sessionId).Err()
	return err == nil
}

func (s *Service) DeleteSession(ctx context.Context, sessionId string) bool {
	err := s.cache.Del(ctx, sessionId).Err()
	return err == nil
}
