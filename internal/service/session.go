package service

import (
	"context"
	"smply/internal/storage"
	"smply/utils"
	"time"
)

func CreateSession(ctx context.Context, identity string) (string, error) {
	sessionId, err := utils.GenerateToken()
	if err != nil {
		return "", err
	}

	err = storage.Cache.Set(ctx, sessionId, identity, 1*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func ValidateSession(ctx context.Context, sessionId string) bool {
	err := storage.Cache.Get(ctx, sessionId).Err()
	return err == nil
}

func DeleteSession(ctx context.Context, sessionId string) bool {
	err := storage.Cache.Del(ctx, sessionId).Err()
	return err == nil
}
