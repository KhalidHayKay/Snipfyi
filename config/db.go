package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, Env.DbUrl)
	if err != nil {
		return err
	}

	err = db.Ping(ctx)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			original TEXT UNIQUE NOT NULL,
			short TEXT UNIQUE NOT NULL,
			visited INTEGER DEFAULT 0,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_visited TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	DB = db
	return nil
}
