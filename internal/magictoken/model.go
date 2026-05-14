package magictoken

import "time"

type MagicToken struct {
	Id        int64      `json:"id"`
	Email     string     `json:"email"`
	TokenHash string     `json:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
