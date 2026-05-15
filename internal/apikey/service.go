package apikey

import (
	"context"
	"log"
	"smply/internal/magictoken"
	"smply/internal/queue"
	"smply/utils"
)

type Service struct {
	repo              Repository
	magicTokenService *magictoken.Service
}

func NewService(repo Repository, magicTokenService *magictoken.Service) *Service {
	return &Service{
		repo,
		magicTokenService,
	}
}

func (s *Service) RequestNew(ctx context.Context, email string) error {
	token, err := s.magicTokenService.Create(ctx, email)
	if err != nil {
		log.Printf("error creating magic token: %v", err)
		return err
	}

	err = queue.EnqueueAPIKeyMagicLinkEmail(ctx, email, token)
	if err != nil {
		log.Printf("failed to enqueue email: %v", err)
	}

	return nil
}

func (s *Service) Create(ctx context.Context, token string) (string, error) {
	magicToken, err := s.magicTokenService.Validate(ctx, token)
	if err != nil {
		return "", err
	}

	err = s.repo.RevokeAll(ctx, magicToken.Email)
	if err != nil {
		log.Printf("error revoking old API keys: %v", err)
		return "", err
	}

	rawkey, err := utils.GenerateAPIKey()
	if err != nil {
		log.Printf("error generating API key: %v", err)
		return "", err
	}

	_, err = s.repo.Create(ctx, magicToken.Email, rawkey)
	if err != nil {
		log.Printf("error creating API key: %v", err)
		return "", err
	}

	return rawkey, nil
}

func (s *Service) Validate(ctx context.Context, key string) (bool, error) {
	keyHash := utils.Hash(key)

	apiKey, err := s.repo.FindByHash(ctx, keyHash)
	if err != nil {
		log.Printf("error finding API key: %v", err)
		return false, err
	}
	if apiKey == nil {
		return false, nil
	}

	if apiKey.RevokedAt != nil {
		return false, nil
	}

	return true, nil
}
