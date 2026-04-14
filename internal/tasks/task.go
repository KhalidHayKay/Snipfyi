package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeAPIKeyMagicLinkEmail = "email:api-key-magic-link"
)

type APIKeyMagicLinkEmailPayload struct {
	Email string
	Token string
}

func NewAPIKeyMagicLinkEmailTask(email string, token string) (*asynq.Task, error) {
	payload, err := json.Marshal(APIKeyMagicLinkEmailPayload{
		Email: email,
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeAPIKeyMagicLinkEmail, payload), nil
}
