package storage

import (
	"context"
	"smply/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func InitCache() (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cache := redis.NewClient(&redis.Options{
		Addr:     config.Env.Redis.Url,
		Password: config.Env.Redis.Password,
		DB:       0,
	})

	_, err := cache.Ping(ctx).Result()
	if err != nil {
		cache.Close()
		return nil, err
	}

	// Todo: remove global variable
	Cache = cache

	return cache, nil
}
