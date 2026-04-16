package service

import (
	"context"
	"log"
	"smply/internal/storage"
	"smply/model"
	"smply/utils"
)

func CreateApiKey(ctx context.Context, token string) (string, error) {
	magicToken, err := ValidateMagicToken(ctx, token)
	if err != nil {
		log.Printf("error validating magic token: %v", err)
		return "", err
	}

	// Revoke existing API keys for this email
	_, err = storage.DB.Exec(ctx, `
		DELETE FROM api_keys
		WHERE owner_email = $1
	`, magicToken.Email)
	if err != nil {
		log.Printf("error deleting existing API key: %v", err)
		return "", err
	}

	key, err := utils.GenerateAPIKey()
	if err != nil {
		return "", err
	}

	var apiKey model.APIKey

	err = storage.DB.QueryRow(ctx, `
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
		log.Printf("error creating API key: %v", err)
		return "", err
	}

	return key, nil
}

func ValidateAPIKey(ctx context.Context, key string) (bool, error) {
	var exists bool

	err := storage.DB.QueryRow(ctx, `
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
