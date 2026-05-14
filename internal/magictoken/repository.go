package magictoken

import "context"

type Repository interface {
	Create(ctx context.Context, email, tokenHash string) error
	GetValid(ctx context.Context, tokenHash string) (*MagicToken, error)
	MarkUsed(ctx context.Context, tokenHash string) error
	MarkAllUsed(ctx context.Context, email string) error
}
