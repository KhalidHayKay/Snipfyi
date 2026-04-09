package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func InitCache() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:     Env.Redis.Url,
		Password: Env.Redis.Password,
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		rdb.Close()
		return err
	}

	Cache = rdb
	return nil
}
