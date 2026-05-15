package app

import (
	"smply/app/middleware"
	"smply/config"
	"smply/internal/admin"
	"smply/internal/apikey"
	"smply/internal/limiter"
	"smply/internal/magictoken"
	"smply/internal/mail"
	"smply/internal/session"
	"smply/internal/stat"
	"smply/internal/url"
	"smply/internal/worker"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

type Handlers struct {
	URL    *url.Handler
	APIKey *apikey.Handler
	Admin  *admin.Handler
	Stat   *stat.Handler
	Worker *worker.Handler
}

type App struct {
	Handlers   Handlers
	Middleware *middleware.Middleware
}

func Bootstrap(db *pgxpool.Pool, cache *redis.Client) App {
	magicTokenRepo := magictoken.NewPostgresRepo(db)
	magicTokenService := magictoken.NewService(magicTokenRepo)

	apikeyRepo := apikey.NewPostgresRepo(db)
	apiKeyService := apikey.NewService(apikeyRepo, magicTokenService)

	sessionService := session.NewService(cache)

	adminService := admin.NewService(magicTokenService, sessionService)

	urlRepo := url.NewPostgresRepo(db)
	urlService := url.NewService(urlRepo)

	statRepo := stat.NewPostresRepo(db)
	statService := stat.NewService(statRepo)

	mailService := mail.NewService(gomail.NewDialer(
		config.Env.Mailer.Host,
		config.Env.Mailer.Port,
		config.Env.Mailer.User,
		config.Env.Mailer.Pass,
	))

	rateLimiterService := limiter.NewService(cache)

	workerHandler := worker.NewHandler(statService, mailService)

	return App{
		Handlers: Handlers{
			APIKey: apikey.NewHandler(apiKeyService),
			Admin:  admin.NewHandler(adminService),
			URL:    url.NewHandler(urlService),
			Stat:   stat.NewHandler(statService),
			Worker: workerHandler,
		},
		Middleware: middleware.NewMiddleware(
			apiKeyService,
			sessionService,
			rateLimiterService,
		),
	}
}
