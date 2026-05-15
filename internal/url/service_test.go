package url

import (
	"context"
	"errors"
	"testing"

	"smply/config"
)

type mockRepository struct {
	getExactFunc   func(ctx context.Context, originalUrl, alias string) (*Url, error)
	storeFunc      func(ctx context.Context, url, alias string) (Url, error)
	getByAliasFunc func(ctx context.Context, alias string) (Url, error)
	calls          []string
}

func (m *mockRepository) GetExact(ctx context.Context, originalUrl, alias string) (*Url, error) {
	m.calls = append(m.calls, "GetExact")
	return m.getExactFunc(ctx, originalUrl, alias)
}

func (m *mockRepository) Store(ctx context.Context, url, alias string) (Url, error) {
	m.calls = append(m.calls, "Store")
	return m.storeFunc(ctx, url, alias)
}

func (m *mockRepository) GetByAlias(ctx context.Context, alias string) (Url, error) {
	m.calls = append(m.calls, "GetByAlias")
	return m.getByAliasFunc(ctx, alias)
}

func ensureUrlConfig() {
	if config.Env == nil {
		config.Env = &config.EnvType{App: config.AppConfig{Url: "http://example.com"}}
	}
	if config.Env.App.Url == "" {
		config.Env.App.Url = "http://example.com"
	}
}

func TestStoreReturnsExistingUrl(t *testing.T) {
	ensureUrlConfig()
	ctx := context.Background()
	repo := &mockRepository{
		getExactFunc: func(ctx context.Context, originalUrl, alias string) (*Url, error) {
			return &Url{Original: originalUrl, Alias: alias}, nil
		},
		storeFunc: func(ctx context.Context, url, alias string) (Url, error) {
			t.Fatal("Store should not be called when an existing URL is returned")
			return Url{}, nil
		},
	}

	svc := NewService(repo)
	got, err := svc.Store(ctx, "https://example.com", "short")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected url, got nil")
	}
	if got.Alias != "short" {
		t.Fatalf("expected alias short, got %q", got.Alias)
	}
	if got.ShortUrl == "" || got.StatUrl == "" {
		t.Fatalf("expected built urls, got short=%q stat=%q", got.ShortUrl, got.StatUrl)
	}
	if len(repo.calls) != 1 || repo.calls[0] != "GetExact" {
		t.Fatalf("expected only GetExact call, got %v", repo.calls)
	}
}

func TestStoreCreatesUrlWhenNotExisting(t *testing.T) {
	ensureUrlConfig()
	ctx := context.Background()
	repo := &mockRepository{
		getExactFunc: func(ctx context.Context, originalUrl, alias string) (*Url, error) {
			return nil, errors.New("not found")
		},
		storeFunc: func(ctx context.Context, url, alias string) (Url, error) {
			return Url{Original: url, Alias: "generated"}, nil
		},
	}

	svc := NewService(repo)
	got, err := svc.Store(ctx, "https://example.com", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Alias == "" {
		t.Fatal("expected generated alias when input alias is empty")
	}
	if got.ShortUrl == "" || got.StatUrl == "" {
		t.Fatal("expected built urls")
	}
	if len(repo.calls) != 2 || repo.calls[0] != "GetExact" || repo.calls[1] != "Store" {
		t.Fatalf("expected GetExact then Store, got %v", repo.calls)
	}
}

func TestStoreReturnsErrorWhenStoreFails(t *testing.T) {
	ensureUrlConfig()
	ctx := context.Background()
	repo := &mockRepository{
		getExactFunc: func(ctx context.Context, originalUrl, alias string) (*Url, error) {
			return nil, errors.New("not found")
		},
		storeFunc: func(ctx context.Context, url, alias string) (Url, error) {
			return Url{}, errors.New("store failed")
		},
	}

	svc := NewService(repo)
	_, err := svc.Store(ctx, "https://example.com", "short")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetByAliasBuildsUrls(t *testing.T) {
	ensureUrlConfig()
	ctx := context.Background()
	repo := &mockRepository{
		getByAliasFunc: func(ctx context.Context, alias string) (Url, error) {
			return Url{Original: "https://example.com", Alias: alias}, nil
		},
	}

	svc := NewService(repo)
	got, err := svc.GetByAlias(ctx, "short")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ShortUrl == "" || got.StatUrl == "" {
		t.Fatal("expected built urls")
	}
	if len(repo.calls) != 1 || repo.calls[0] != "GetByAlias" {
		t.Fatalf("expected GetByAlias call, got %v", repo.calls)
	}
}

func TestGetByAliasReturnsErrorWhenRepoFails(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepository{
		getByAliasFunc: func(ctx context.Context, alias string) (Url, error) {
			return Url{}, errors.New("repo failure")
		},
	}

	svc := NewService(repo)
	_, err := svc.GetByAlias(ctx, "short")
	if err == nil {
		t.Fatal("expected error")
	}
}
