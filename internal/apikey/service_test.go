package apikey

import (
	"context"
	"errors"
	"testing"
	"time"

	"smply/internal/magictoken"
	"smply/utils"
)

type mockRepo struct {
	revokeAllFunc  func(ctx context.Context, email string) error
	createFunc     func(ctx context.Context, email, key string) (APIKey, error)
	findByHashFunc func(ctx context.Context, keyHash string) (*APIKey, error)
	calls          []string
}

func (m *mockRepo) RevokeAll(ctx context.Context, email string) error {
	m.calls = append(m.calls, "RevokeAll")
	return m.revokeAllFunc(ctx, email)
}

func (m *mockRepo) Create(ctx context.Context, email, key string) (APIKey, error) {
	m.calls = append(m.calls, "Create")
	return m.createFunc(ctx, email, key)
}

func (m *mockRepo) FindByHash(ctx context.Context, keyHash string) (*APIKey, error) {
	m.calls = append(m.calls, "FindByHash")
	return m.findByHashFunc(ctx, keyHash)
}

type mockMagicTokenRepo struct {
	getValidFunc func(ctx context.Context, tokenHash string) (*magictoken.MagicToken, error)
	markUsedFunc func(ctx context.Context, tokenHash string) error
}

func (m *mockMagicTokenRepo) Create(ctx context.Context, email, tokenHash string) error {
	return nil
}

func (m *mockMagicTokenRepo) GetValid(ctx context.Context, tokenHash string) (*magictoken.MagicToken, error) {
	return m.getValidFunc(ctx, tokenHash)
}

func (m *mockMagicTokenRepo) MarkUsed(ctx context.Context, tokenHash string) error {
	return m.markUsedFunc(ctx, tokenHash)
}

func (m *mockMagicTokenRepo) MarkAllUsed(ctx context.Context, email string) error {
	return nil
}

func TestCreateRevokesAllBeforeCreate(t *testing.T) {
	ctx := context.Background()
	magicRepo := &mockMagicTokenRepo{
		getValidFunc: func(ctx context.Context, tokenHash string) (*magictoken.MagicToken, error) {
			return &magictoken.MagicToken{Email: "alice@example.com"}, nil
		},
		markUsedFunc: func(ctx context.Context, tokenHash string) error { return nil },
	}
	magicService := magictoken.NewService(magicRepo)

	repo := &mockRepo{
		revokeAllFunc: func(ctx context.Context, email string) error {
			if email != "alice@example.com" {
				t.Fatalf("unexpected email %q", email)
			}
			return nil
		},
		createFunc: func(ctx context.Context, email, key string) (APIKey, error) {
			if len(key) == 0 {
				t.Fatal("expected generated key")
			}
			return APIKey{OwnerEmail: email, KeyHash: utils.Hash(key)}, nil
		},
	}

	svc := NewService(repo, magicService)
	apiKey, err := svc.Create(ctx, "token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if apiKey == "" {
		t.Fatal("expected key to be returned")
	}
	if len(repo.calls) != 2 || repo.calls[0] != "RevokeAll" || repo.calls[1] != "Create" {
		t.Fatalf("expected revoke then create, got %v", repo.calls)
	}
}

func TestCreateReturnsErrorWhenValidateFails(t *testing.T) {
	ctx := context.Background()
	magicRepo := &mockMagicTokenRepo{
		getValidFunc: func(ctx context.Context, tokenHash string) (*magictoken.MagicToken, error) {
			return nil, errors.New("invalid token")
		},
		markUsedFunc: func(ctx context.Context, tokenHash string) error { return nil },
	}
	magicService := magictoken.NewService(magicRepo)
	repo := &mockRepo{}

	svc := NewService(repo, magicService)
	_, err := svc.Create(ctx, "token")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateReturnsErrorWhenRevokeFails(t *testing.T) {
	ctx := context.Background()
	magicRepo := &mockMagicTokenRepo{
		getValidFunc: func(ctx context.Context, tokenHash string) (*magictoken.MagicToken, error) {
			return &magictoken.MagicToken{Email: "alice@example.com"}, nil
		},
		markUsedFunc: func(ctx context.Context, tokenHash string) error { return nil },
	}
	magicService := magictoken.NewService(magicRepo)
	repo := &mockRepo{
		revokeAllFunc: func(ctx context.Context, email string) error {
			return errors.New("revoke failed")
		},
		createFunc: func(ctx context.Context, email, key string) (APIKey, error) {
			t.Fatal("Create should not be called when revoke fails")
			return APIKey{}, nil
		},
	}

	svc := NewService(repo, magicService)
	_, err := svc.Create(ctx, "token")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateReturnsTrueWhenKeyFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{
		findByHashFunc: func(ctx context.Context, keyHash string) (*APIKey, error) {
			return &APIKey{OwnerEmail: "alice@example.com", KeyHash: keyHash}, nil
		},
	}

	svc := NewService(repo, nil)
	ok, err := svc.Validate(ctx, "key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key validation to succeed")
	}
}

func TestValidateReturnsFalseWhenKeyNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{
		findByHashFunc: func(ctx context.Context, keyHash string) (*APIKey, error) {
			return nil, nil
		},
	}

	svc := NewService(repo, nil)
	ok, err := svc.Validate(ctx, "key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected key validation to fail for missing key")
	}
}

func TestValidateRejectsRevokedKey(t *testing.T) {
	ctx := context.Background()
	revokedAt := time.Now()
	repo := &mockRepo{
		findByHashFunc: func(ctx context.Context, keyHash string) (*APIKey, error) {
			return &APIKey{OwnerEmail: "alice@example.com", KeyHash: keyHash, RevokedAt: &revokedAt}, nil
		},
	}

	svc := NewService(repo, nil)
	ok, err := svc.Validate(ctx, "key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected revoked key to be rejected")
	}
}

func TestValidateReturnsErrorWhenRepoFails(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{
		findByHashFunc: func(ctx context.Context, keyHash string) (*APIKey, error) {
			return nil, errors.New("find failed")
		},
	}

	svc := NewService(repo, nil)
	_, err := svc.Validate(ctx, "key")
	if err == nil {
		t.Fatal("expected error")
	}
}
