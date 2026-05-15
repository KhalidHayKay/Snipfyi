package url

import "context"

type Repository interface {
	Store(ctx context.Context, url string, alias string) (Url, error)
	GetExact(ctx context.Context, originalUrl, alias string) (*Url, error)
	GetByAlias(ctx context.Context, alias string) (Url, error)
}
