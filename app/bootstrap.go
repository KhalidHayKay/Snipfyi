package app

import (
	"smply/internal/apikey"
	"smply/internal/magictoken"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	APIKey *apikey.Handler
}

type Services struct {
	APIKey *apikey.Service
}

func Bootstrap(db *pgxpool.Pool) (Handlers, Services) {
	magicTokenRepo := magictoken.NewPostgresRepo(db)
	magicTokenService := magictoken.NewService(magicTokenRepo)

	apikeyRepo := apikey.NewPostgresRepo(db)
	apiKeyService := apikey.NewService(apikeyRepo, *magicTokenService)

	return Handlers{
			APIKey: apikey.NewHandler(*apiKeyService),
		}, Services{
			APIKey: apiKeyService,
		}
}
