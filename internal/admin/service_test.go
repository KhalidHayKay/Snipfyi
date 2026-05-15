package admin

import (
	"context"
	"errors"
	"testing"

	"smply/config"
	"smply/internal/magictoken"
	"smply/internal/session"

	"github.com/redis/go-redis/v9"
)

type mockMagicTokenRepo struct {
	createFunc   func(ctx context.Context, email, tokenHash string) error
	getValidFunc func(ctx context.Context, tokenHash string) (*magictoken.MagicToken, error)
	markUsedFunc func(ctx context.Context, tokenHash string) error
}

func (m *mockMagicTokenRepo) Create(ctx context.Context, email, tokenHash string) error {
	return m.createFunc(ctx, email, tokenHash)
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

func ensureAdminConfig() {
	if config.Env == nil {
		config.Env = &config.EnvType{App: config.AppConfig{Url: "http://example.com"}}
	}
	if config.Env.App.Url == "" {
		config.Env.App.Url = "http://example.com"
	}
}

func TestAttemptMagicLinkLoginRejectsNonAdminEmail(t *testing.T) {
	ensureAdminConfig()
	oldAdminEmail := config.Env.AdminEmail
	config.Env.AdminEmail = "admin@example.com"
	defer func() { config.Env.AdminEmail = oldAdminEmail }()

	svc := NewService(magictoken.NewService(&mockMagicTokenRepo{}), nil)
	err := svc.AttemptMagicLinkLogin(context.Background(), "user@example.com")
	if err == nil {
		t.Fatal("expected error for non-admin email")
	}
}

func TestAttemptMagicLinkLoginReturnsInternalServerErrorOnCreateFailure(t *testing.T) {
	ensureAdminConfig()
	oldAdminEmail := config.Env.AdminEmail
	config.Env.AdminEmail = "admin@example.com"
	defer func() { config.Env.AdminEmail = oldAdminEmail }()

	magicRepo := &mockMagicTokenRepo{
		createFunc: func(ctx context.Context, email, tokenHash string) error {
			return errors.New("create failed")
		},
	}
	svc := NewService(magictoken.NewService(magicRepo), nil)
	err := svc.AttemptMagicLinkLogin(context.Background(), "admin@example.com")
	if err == nil {
		t.Fatal("expected internal server error")
	}
}

func TestAuthenticateReturnsGenericErrorWhenValidateFails(t *testing.T) {
	magicRepo := &mockMagicTokenRepo{
		getValidFunc: func(ctx context.Context, tokenHash string) (*magictoken.MagicToken, error) {
			return nil, errors.New("invalid token")
		},
		markUsedFunc: func(ctx context.Context, tokenHash string) error { return nil },
	}
	svc := NewService(magictoken.NewService(magicRepo), nil)
	_, err := svc.Authenticate(context.Background(), "token")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAuthenticateReturnsInternalServerErrorWhenSessionCreateFails(t *testing.T) {
	magicRepo := &mockMagicTokenRepo{
		getValidFunc: func(ctx context.Context, tokenHash string) (*magictoken.MagicToken, error) {
			return &magictoken.MagicToken{Email: "admin@example.com"}, nil
		},
		markUsedFunc: func(ctx context.Context, tokenHash string) error { return nil },
	}
	cache := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	svc := NewService(magictoken.NewService(magicRepo), session.NewService(cache))
	_, err := svc.Authenticate(context.Background(), "token")
	if err == nil {
		t.Fatal("expected error")
	}
}
