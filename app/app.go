package app

import (
	"log"
	"net/http"
	"smply/config"
	"smply/internal/queue"
	"smply/internal/tasks"
	"smply/internal/worker"

	"github.com/hibiken/asynq"
)

func Start() {
	config.LoadEnv()

	if err := config.InitDB(); err != nil {
		log.Fatal(err)
	}

	if err := config.InitCache(); err != nil {
		log.Fatal(err)
	}

	queue.Init()

	router := setupRouter()

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":"+config.Env.App.Port, router))
}

func StartWorker() {
	config.LoadEnv()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.Env.Redis.Url,
			Password: config.Env.Redis.Password,
			DB:       1,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 10,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeAPIKeyMagicLinkEmail, worker.HandleAPIKeyMagicLinkEmail)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
