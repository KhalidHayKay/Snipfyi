package worker

import (
	"context"
	"encoding/json"
	"log"
	"smply/internal/mail"
	"smply/internal/stat"
	"smply/internal/tasks"

	"github.com/hibiken/asynq"
)

type Handler struct {
	statService *stat.Service
	mailService *mail.Service
}

func NewHandler(statService *stat.Service, mailService *mail.Service) *Handler {
	return &Handler{statService, mailService}
}

func (h *Handler) APIKeyMagicLinkEmail(ctx context.Context, task *asynq.Task) error {
	var payload tasks.APIKeyMagicLinkEmailPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		log.Printf("Failed to unmarshal magic link emails task payload: %v", err)
		return err
	}

	return h.mailService.SendAPIKeyMagicLink(payload.Email, payload.Token)
}

func (h *Handler) StatsUpdate(ctx context.Context, task *asynq.Task) error {
	var payload tasks.StatsUpdatePayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		log.Printf("Failed to unmarshal stats update task payload: %v", err)
		return err
	}

	return h.statService.Run(ctx,
		payload.UrlAlias,
		payload.Referer,
		payload.UserAgent,
		payload.Timestamp,
	)
}

func (h *Handler) AdminLoginMagicLinkEmail(ctx context.Context, task *asynq.Task) error {
	var payload tasks.AdminLoginMagicLinkEmailPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		log.Printf("Failed to unmarshal admin login magic link email task payload: %v", err)
		return err
	}

	return h.mailService.SendAdminLoginMagicLink(payload.Email, payload.Token)
}
