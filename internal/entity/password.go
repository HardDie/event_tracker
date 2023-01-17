package entity

import "time"

type Password struct {
	ID             int32      `json:"id"`
	UserID         int32      `json:"userId"`
	PasswordHash   string     `json:"passwordHash"`
	FailedAttempts int32      `json:"failedAttempts"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
	BlockedAt      *time.Time `json:"blockedAt"`
}
