package url

import (
	"context"
	"log"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Store(ctx context.Context, longUrl, alias string) (*Url, error) {
	existingUrl, _ := s.repo.GetExact(ctx, longUrl, alias)
	if existingUrl != nil {
		existingUrl.BuildUrls()
		return existingUrl, nil
	}

	url, err := s.repo.Store(ctx, longUrl, alias)
	if err != nil {
		log.Printf("Error storing url: %v", err)
		return nil, err
	}

	url.BuildUrls()
	return &url, nil
}

func (s *Service) GetByAlias(ctx context.Context, alias string) (*Url, error) {
	url, err := s.repo.GetByAlias(ctx, alias)
	if err != nil {
		log.Printf("Error fetching url by alias: %v", err)
		return nil, err
	}

	url.BuildUrls()
	return &url, nil
}
