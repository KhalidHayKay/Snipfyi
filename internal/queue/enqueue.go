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
