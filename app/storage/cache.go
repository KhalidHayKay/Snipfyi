package storage

import (
	"context"
	"smply/config"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redis := redis.NewClient(&redis.Options{
		Addr:     config.Env.Redis.Url,
		Password: config.Env.Redis.Password,
		DB:       0,
	})

	_, err := redis.Ping(ctx).Result()
	if err != nil {
		redis.Close()
		return nil, err
	}

	return redis, nil
}
