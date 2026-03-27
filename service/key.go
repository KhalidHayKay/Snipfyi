package service

import (
	"context"
	"smply/config"
	"smply/model"
	"smply/utils"
)

func RequestApiKey(ctx context.Context, email string) error {
	token, err := utils.GenerateMagicToken()
	if err != nil {
		return err
	}

	var magicToken model.MagicToken

	err = config.DB.QueryRow(ctx, `
		INSERT INTO magic_tokens (email, token_hash, expires_at)
		VALUES ($1, $2, NOW() + INTERVAL '15 minutes')
		RETURNING id, email, token_hash, expires_at
	`, email, utils.Hash(token)).Scan(
		&magicToken.Id,
		&magicToken.Email,
		&magicToken.TokenHash,
		&magicToken.ExpiresAt,
	)

	go func() {
		sendMagicLinkEmail(email, token)
	}()

	if err != nil {
		return err
	}

	return nil
}

func CreateApiKey(ctx context.Context, token string) (string, error) {
	tokenHash := utils.Hash(token)

	var magicToken model.MagicToken

	err := config.DB.QueryRow(ctx, `
		SELECT id, email, token_hash, expires_at
		FROM magic_tokens
		WHERE token_hash = $1 AND expires_at > NOW()
	`, tokenHash).Scan(
		&magicToken.Id,
		&magicToken.Email,
		&magicToken.TokenHash,
		&magicToken.ExpiresAt,
	)
	if err != nil {
		return "", err
	}

	key, err := utils.GenerateAPIKey()
	if err != nil {
		return "", err
	}

	var apiKey model.APIKey

	err = config.DB.QueryRow(ctx, `
		INSERT INTO api_keys (owner_email, key_hash, created_at, expires_at)
		VALUES ($1, $2, NOW(), NOW() + INTERVAL '30 days')
		RETURNING id, owner_email, key_hash, created_at, expires_at
	`, magicToken.Email, utils.Hash(key)).Scan(
		&apiKey.Id,
		&apiKey.OwnerEmail,
		&apiKey.KeyHash,
		&apiKey.CreatedAt,
		&apiKey.ExpiresAt,
	)

	if err != nil {
		return "", err
	}

	return key, nil
}
