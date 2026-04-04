package model

import "time"

type APIKey struct {
	Id         int64      `json:"id"`
	OwnerEmail string     `json:"owner_email"`
	KeyHash    string     `json:"key_hash"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}

type MagicToken struct {
	Id        int64      `json:"id"`
	Email     string     `json:"email"`
	TokenHash string     `json:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
