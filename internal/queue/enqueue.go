package queue

import (
	"context"
	"smply/internal/tasks"
	"time"

	"github.com/hibiken/asynq"
)

func EnqueueAPIKeyMagicLinkEmail(ctx context.Context, email, token string) error {
	t, err := tasks.NewAPIKeyMagicLinkEmailTask(email, token)
	if err != nil {
		return err
	}

	_, err = client.EnqueueContext(ctx, t,
		asynq.Queue("critical"),
		asynq.MaxRetry(5),
		asynq.Timeout(30*time.Second),
	)
	return err
}

func EnqueueStatsUpdate(ctx context.Context, urlAlias, referer, userAgent, ipAddress string, timestamp time.Time) error {
	t, err := tasks.NewStatsUpdateTask(urlAlias, referer, userAgent, ipAddress, timestamp)
	if err != nil {
		return err
	}

	_, err = client.EnqueueContext(ctx, t,
		asynq.Queue("default"),
		asynq.MaxRetry(3),
		asynq.Timeout(10*time.Second),
	)
	return err
}

func EnqueueAdminLoginMagicLinkEmail(ctx context.Context, email, token string) error {
	t, err := tasks.NewAdminLoginMagicLinkEmailTask(email, token)
	if err != nil {
		return err
	}

	_, err = client.EnqueueContext(ctx, t,
		asynq.Queue("critical"),
		asynq.MaxRetry(5),
		asynq.Timeout(30*time.Second),
	)
	return err
}
