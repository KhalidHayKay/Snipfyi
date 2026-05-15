package app

import (
	"log"
	"net/http"
	"smply/app/storage"
	"smply/config"
	"smply/internal/queue"
	"smply/internal/tasks"

	"github.com/hibiken/asynq"
)

func Start() {
	config.LoadEnv()

	pgsql, err := storage.InitPostgres()
	if err != nil {
		log.Fatal(err)
	}

	redis, err := storage.InitRedis()
	if err != nil {
		log.Fatal(err)
	}

	queue.Init()

	app := Bootstrap(pgsql, redis)

	router := setupRouter(app.Handlers, app.Middleware)

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":"+config.Env.App.Port, router))
}

func StartWorker() {
	config.LoadEnv()

	pgsql, err := storage.InitPostgres()
	if err != nil {
		log.Fatal(err)
	}

	redis, err := storage.InitRedis()
	if err != nil {
		log.Fatal(err)
	}

	app := Bootstrap(pgsql, redis)
	worker := app.Handlers.Worker

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
				"default":  5,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeAPIKeyMagicLinkEmail, worker.APIKeyMagicLinkEmail)
	mux.HandleFunc(tasks.TypeStatsUpdate, worker.StatsUpdate)
	mux.HandleFunc(tasks.TypeAdminLoginMagicLinkEmail, worker.AdminLoginMagicLinkEmail)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
