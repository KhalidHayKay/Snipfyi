package worker

import (
	"context"
	"encoding/json"
	"log"
	"smply/internal/service"
	"smply/internal/tasks"

	"github.com/hibiken/asynq"
)

func HandleAPIKeyMagicLinkEmail(ctx context.Context, task *asynq.Task) error {
	var payload tasks.APIKeyMagicLinkEmailPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		log.Printf("Failed to unmarshal magic link emails task payload: %v", err)
		return err
	}

	return service.SendMagicLinkEmail(payload.Email, payload.Token)
}

func HandleStatsUpdate(ctx context.Context, task *asynq.Task) error {
	var payload tasks.StatsUpdatePayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		log.Printf("Failed to unmarshal stats update task payload: %v", err)
		return err
	}

	return service.RunStats(ctx,
		payload.UrlAlias,
		payload.Referrer,
		payload.UserAgent,
		payload.IpAddress,
		payload.Timestamp,
	)
}
