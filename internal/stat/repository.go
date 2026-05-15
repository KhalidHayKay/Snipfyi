package stat

import (
	"context"
	"time"
)

type Repository interface {
	Run(ctx context.Context, alias, referer, userAgent string, timestamp time.Time) error
	Get(ctx context.Context, alias string) (Stats, error)
	GetAdmin(ctx context.Context) (AdminStats, error)
}
