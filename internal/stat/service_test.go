package stat

import (
	"context"
	"errors"
	"testing"
	"time"

	"smply/config"
)

type mockRepository struct {
	runFunc      func(ctx context.Context, alias, referer, userAgent string, timestamp time.Time) error
	getFunc      func(ctx context.Context, alias string) (Stats, error)
	getAdminFunc func(ctx context.Context) (AdminStats, error)
	calls        []string
}

func (m *mockRepository) Run(ctx context.Context, alias, referer, userAgent string, timestamp time.Time) error {
	m.calls = append(m.calls, "Run")
	return m.runFunc(ctx, alias, referer, userAgent, timestamp)
}

func (m *mockRepository) Get(ctx context.Context, alias string) (Stats, error) {
	m.calls = append(m.calls, "Get")
	return m.getFunc(ctx, alias)
}

func (m *mockRepository) GetAdmin(ctx context.Context) (AdminStats, error) {
	m.calls = append(m.calls, "GetAdmin")
	return m.getAdminFunc(ctx)
}

func TestRunReturnsErrorOnRepoFailure(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepository{
		runFunc: func(ctx context.Context, alias, referer, userAgent string, timestamp time.Time) error {
			return errors.New("run failed")
		},
	}

	svc := NewService(repo)
	err := svc.Run(ctx, "short", "ref", "agent", time.Now())
	if err == nil {
		t.Fatal("expected error")
	}
}

func ensureStatConfig() {
	if config.Env == nil {
		config.Env = &config.EnvType{App: config.AppConfig{Url: "http://example.com"}}
	}
	if config.Env.App.Url == "" {
		config.Env.App.Url = "http://example.com"
	}
}

func TestGetBuildsShortUrl(t *testing.T) {
	ensureStatConfig()
	ctx := context.Background()
	repo := &mockRepository{
		getFunc: func(ctx context.Context, alias string) (Stats, error) {
			return Stats{Alias: alias}, nil
		},
	}

	svc := NewService(repo)
	got, err := svc.Get(ctx, "short")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ShortUrl == "" {
		t.Fatal("expected short url to be built")
	}
	if len(repo.calls) != 1 || repo.calls[0] != "Get" {
		t.Fatalf("expected Get call, got %v", repo.calls)
	}
}

func TestGetReturnsErrorWhenRepoFails(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepository{
		getFunc: func(ctx context.Context, alias string) (Stats, error) {
			return Stats{}, errors.New("get failed")
		},
	}

	svc := NewService(repo)
	_, err := svc.Get(ctx, "short")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetAdminSuccess(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepository{
		getAdminFunc: func(ctx context.Context) (AdminStats, error) {
			return AdminStats{TotalUrls: 1}, nil
		},
	}

	svc := NewService(repo)
	got, err := svc.GetAdmin(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TotalUrls != 1 {
		t.Fatalf("expected TotalUrls 1, got %d", got.TotalUrls)
	}
}

func TestGetAdminReturnsErrorWhenRepoFails(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepository{
		getAdminFunc: func(ctx context.Context) (AdminStats, error) {
			return AdminStats{}, errors.New("get admin failed")
		},
	}

	svc := NewService(repo)
	_, err := svc.GetAdmin(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
}
