package apikey

import (
	"context"
	"log"
	"smply/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db}
}

func (r PostgresRepo) Create(ctx context.Context, email, key string) (APIKey, error) {
	var apiKey APIKey

	err := r.db.QueryRow(ctx, `
		INSERT INTO api_keys (owner_email, key_hash, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, owner_email, key_hash, created_at
	`, email, utils.Hash(key)).Scan(
		&apiKey.Id,
		&apiKey.OwnerEmail,
		&apiKey.KeyHash,
		&apiKey.CreatedAt,
	)

	if err != nil {
		log.Printf("error creating API key: %v", err)
		return APIKey{}, err
	}

	return apiKey, nil
}

func (r PostgresRepo) FindByHash(ctx context.Context, keyHash string) (*APIKey, error) {
	var apiKey APIKey

	err := r.db.QueryRow(ctx, `
		SELECT id, owner_email, key_hash, created_at
		FROM api_keys
		WHERE key_hash = $1
	`, keyHash).Scan(
		&apiKey.Id,
		&apiKey.OwnerEmail,
		&apiKey.KeyHash,
		&apiKey.CreatedAt,
	)

	if err != nil {
		log.Printf("error finding API key: %v", err)
		return nil, err
	}

	return &apiKey, nil
}

func (r PostgresRepo) RevokeAll(ctx context.Context, email string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE api_keys
			SET revoked_at = NOW()
			WHERE owner_email = $1
	`, email)
	if err != nil {
		log.Printf("error revoking API keys: %v", err)
		return err
	}

	return nil
}
