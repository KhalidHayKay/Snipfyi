package magictoken

import (
	"context"
	"log"
	"smply/utils"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(ctx context.Context, email string) (string, error) {
	token, err := utils.GenerateToken()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	err = s.repo.Create(ctx, email, utils.Hash(token))
	if err != nil {
		log.Printf("Error creating magic token: %v", err)
		return "", err
	}

	return token, nil
}

func (s *Service) Validate(ctx context.Context, token string) (*MagicToken, error) {
	magicToken, err := s.repo.GetValid(ctx, utils.Hash(token))
	if err != nil {
		log.Printf("Error validating magic token: %v", err)
		return nil, err
	}

	err = s.repo.MarkUsed(ctx, magicToken.TokenHash)
	if err != nil {
		log.Printf("Error marking magic token as used: %v", err)
		return nil, err
	}

	return magicToken, nil
}
