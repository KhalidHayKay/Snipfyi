package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeAPIKeyMagicLinkEmail     = "email:api-key-magic-link"
	TypeAdminLoginMagicLinkEmail = "email:admin-login-magic-link"
	TypeStatsUpdate              = "stats:update"
)

type APIKeyMagicLinkEmailPayload struct {
	Email string
	Token string
}

type AdminLoginMagicLinkEmailPayload struct {
	Email string
	Token string
}

type StatsUpdatePayload struct {
	UrlAlias  string
	Referer   string
	UserAgent string
	IpAddress string
	Timestamp time.Time
}

func NewAPIKeyMagicLinkEmailTask(email, token string) (*asynq.Task, error) {
	payload, err := json.Marshal(APIKeyMagicLinkEmailPayload{
		Email: email,
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeAPIKeyMagicLinkEmail, payload), nil
}

func NewStatsUpdateTask(urlAlias, referer, userAgent, ipAddress string, timestamp time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(StatsUpdatePayload{
		UrlAlias:  urlAlias,
		Referer:   referer,
		UserAgent: userAgent,
		IpAddress: ipAddress,
		Timestamp: timestamp,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeStatsUpdate, payload), nil
}

func NewAdminLoginMagicLinkEmailTask(email string, token string) (*asynq.Task, error) {
	payload, err := json.Marshal(APIKeyMagicLinkEmailPayload{
		Email: email,
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeAdminLoginMagicLinkEmail, payload), nil
}
