package stat

import (
	"context"
	"log"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Run(ctx context.Context, alias, referer, userAgent string, timestamp time.Time) error {
	if err := s.repo.Run(ctx, alias, referer, userAgent, timestamp); err != nil {
		log.Printf("Failed to run stats for alias %s: %v", alias, err)
		return err
	}

	return nil
}

func (s *Service) Get(ctx context.Context, alias string) (*Stats, error) {
	stats, err := s.repo.Get(ctx, alias)
	if err != nil {
		log.Printf("Failed to get stats for alias %s: %v", alias, err)
		return nil, err
	}

	stats.BuildShortUrl()
	log.Println(stats)
	return &stats, nil
}

func (s *Service) GetAdmin(ctx context.Context) (*AdminStats, error) {
	stats, err := s.repo.GetAdmin(ctx)
	if err != nil {
		log.Printf("Failed to get admin stats: %v", err)
		return nil, err
	}

	return &stats, nil
}
