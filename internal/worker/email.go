package worker

import (
	"context"
	"encoding/json"
	"smply/internal/tasks"
	"smply/service"

	"github.com/hibiken/asynq"
)

func HandleAPIKeyMagicLinkEmail(ctx context.Context, task *asynq.Task) error {
	var payload tasks.APIKeyMagicLinkEmailPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	err = service.SendMagicLinkEmail(payload.Email, payload.Token)
	if err != nil {
		return err
	}

	return nil
}
