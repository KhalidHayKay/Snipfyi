package service

import (
	"context"
	"smply/internal/storage"
	"smply/model"
	"smply/utils"
)

func CreateMagicToken(ctx context.Context, email string) (string, error) {
	_, err := storage.DB.Exec(ctx, `
		UPDATE magic_tokens
		SET used_at = NOW()
		WHERE email = $1
		  AND expires_at > NOW()
		  AND used_at IS NULL
	`, email)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken()
	if err != nil {
		return "", err
	}

	var magicToken model.MagicToken

	err = storage.DB.QueryRow(ctx, `
		INSERT INTO magic_tokens (email, token_hash, expires_at, created_at)
		VALUES ($1, $2, NOW() + INTERVAL '15 minutes', NOW())
		RETURNING id, email, token_hash, expires_at
	`, email, utils.Hash(token)).Scan(
		&magicToken.Id,
		&magicToken.Email,
		&magicToken.TokenHash,
		&magicToken.ExpiresAt,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateMagicToken(ctx context.Context, token string) (model.MagicToken, error) {
	tokenHash := utils.Hash(token)

	var magicToken model.MagicToken

	err := storage.DB.QueryRow(ctx, `
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
		return model.MagicToken{}, err
	}

	// Mark magic link Token as used
	_, err = storage.DB.Exec(ctx, `
		UPDATE magic_tokens
		SET used_at = NOW()
		WHERE id = $1
	`, magicToken.Id)
	if err != nil {
		return model.MagicToken{}, err
	}

	return magicToken, nil
}
