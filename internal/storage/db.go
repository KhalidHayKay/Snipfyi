package storage

import (
	"context"
	"smply/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	DB, err = pgxpool.New(ctx, config.Env.DbUrl)
	if err != nil {
		return err
	}

	err = DB.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}
