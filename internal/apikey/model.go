package apikey

import "time"

type APIKey struct {
	Id         int64      `json:"id"`
	OwnerEmail string     `json:"owner_email"`
	KeyHash    string     `json:"key_hash"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}
