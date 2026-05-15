package admin

import (
	"context"
	"errors"
	"log"
	"smply/config"
	"smply/internal/magictoken"
	"smply/internal/queue"
	"smply/internal/session"
)

type Service struct {
	magicTokenService *magictoken.Service
	sessionService    *session.Service
}

func NewService(magicTokenService *magictoken.Service, sessionService *session.Service) *Service {
	return &Service{
		magicTokenService: magicTokenService,
		sessionService:    sessionService,
	}
}

func (s *Service) AttemptMagicLinkLogin(ctx context.Context, email string) error {
	if email != config.Env.AdminEmail {
		return errors.New("Email is not approved")
	}

	token, err := s.magicTokenService.Create(ctx, email)
	if err != nil {
		log.Printf("Error creating magic token for %s: %v", email, err)
		return errors.New("Internal server error")
	}

	err = queue.EnqueueAdminLoginMagicLinkEmail(ctx, email, token)
	if err != nil {
		log.Println("Error enqueuing admin login magic link email:", err)
		return errors.New("Internal server error")
	}

	return nil
}

func (s *Service) Authenticate(ctx context.Context, token string) (string, error) {
	magicToken, err := s.magicTokenService.Validate(ctx, token)
	if err != nil {
		log.Printf("Error validating magic token: %v", err)
		return "", errors.New("Invalid or expired token")
	}

	sessionId, err := s.sessionService.CreateSession(ctx, magicToken.Email)
	if err != nil {
		log.Println("Error creating session:", err)
		return "", errors.New("Internal server error")
	}

	return sessionId, nil
}
