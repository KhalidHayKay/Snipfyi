package service

import (
	"context"
	"smply/config"
	"smply/model"
	"smply/utils"
)

func RequestApiKey(ctx context.Context, email string) error {
	_, err := config.DB.Exec(ctx, `
		UPDATE magic_tokens
		SET used_at = NOW()
		WHERE email = $1
		  AND expires_at > NOW()
		  AND used_at IS NULL
	`, email)
	if err != nil {
		return err
	}

	token, err := utils.GenerateMagicToken()
	if err != nil {
		return err
	}

	var magicToken model.MagicToken

	err = config.DB.QueryRow(ctx, `
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
		return err
	}

	go func() {
		sendMagicLinkEmail(email, token)
	}()

	return nil
}

func CreateApiKey(ctx context.Context, token string) (string, error) {
	tokenHash := utils.Hash(token)

	var magicToken model.MagicToken

	err := config.DB.QueryRow(ctx, `
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
		return "", err
	}

	// Mark magic link Token as used
	_, err = config.DB.Exec(ctx, `
		UPDATE magic_tokens
		SET used_at = NOW()
		WHERE id = $1
	`, magicToken.Id)
	if err != nil {
		return "", err
	}

	// Revoke existing API keys for this email
	_, err = config.DB.Exec(ctx, `
		UPDATE api_keys
		SET expires_at = NOW()
		WHERE owner_email = $1
	`, magicToken.Email)
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

func ValidateAPIKey(ctx context.Context, key string) (bool, error) {
	if key == config.Env.InternalApiKey {
		return true, nil
	}

	var exists bool

	err := config.DB.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM api_keys
			WHERE key_hash = $1 AND expires_at > NOW()
		)
	`, utils.Hash(key)).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
