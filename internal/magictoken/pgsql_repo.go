package magictoken

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pgsql *pgxpool.Pool
}

func NewPostgresRepo(pgsql *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pgsql}
}

func (r PostgresRepo) Create(ctx context.Context, email, tokenHash string) error {
	_, err := r.pgsql.Exec(ctx, `
		INSERT INTO magic_tokens (email, token_hash, expires_at, created_at)
			VALUES ($1, $2, NOW() + INTERVAL '15 minutes', NOW())
	`, email, tokenHash)
	if err != nil {
		return err
	}

	return nil
}

func (r PostgresRepo) GetValid(ctx context.Context, tokenHash string) (*MagicToken, error) {
	var magicToken MagicToken

	err := r.pgsql.QueryRow(ctx, `
		SELECT id, email, token_hash, expires_at
		FROM magic_tokens
		WHERE token_hash = $1 
			AND expires_at > NOW() 
			AND used_at IS NULL
	`, tokenHash).Scan(
		&magicToken.Id,
		&magicToken.Email,
		&magicToken.TokenHash,
		&magicToken.ExpiresAt,
	)
	if err != nil {
		return &MagicToken{}, err
	}

	return &magicToken, nil
}

func (r PostgresRepo) MarkUsed(ctx context.Context, tokenHash string) error {
	_, err := r.pgsql.Exec(ctx, `
		UPDATE magic_tokens
			SET used_at = NOW()
			WHERE token_hash = $1
	`, tokenHash)
	if err != nil {
		return err
	}

	return nil
}

func (r PostgresRepo) MarkAllUsed(ctx context.Context, email string) error {
	_, err := r.pgsql.Exec(ctx, `
		UPDATE magic_tokens
		SET used_at = NOW()
		WHERE email = $1
		  AND expires_at > NOW()
		  AND used_at IS NULL
	`, email)
	if err != nil {
		return err
	}

	return nil
}
