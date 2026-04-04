package config

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

		CREATE TABLE IF NOT EXISTS api_keys (
			id BIGSERIAL PRIMARY KEY,
			owner_email TEXT NOT NULL UNIQUE,
			key_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL,
			last_used_at TIMESTAMP NULL,
			revoked_at TIMESTAMP NULL
		);

		CREATE TABLE IF NOT EXISTS magic_tokens (
			id BIGSERIAL PRIMARY KEY,
			email TEXT NOT NULL,
			token_hash TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			used_at TIMESTAMP NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	DB = db
	return nil
}
