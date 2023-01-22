package entity

import "time"

type Friend struct {
	ID         int32      `json:"id"`
	UserID     int32      `json:"userId"`
	WithUserID int32      `json:"withUserId"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
}
