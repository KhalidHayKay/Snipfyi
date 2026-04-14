package queue

import (
	"smply/config"

	"github.com/hibiken/asynq"
)

var client *asynq.Client

func Init() {
	client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     config.Env.Redis.Url,
		Password: config.Env.Redis.Password,
		DB:       1,
	})
}
