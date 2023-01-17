package entity

import "time"

type Session struct {
	ID          int32      `json:"id"`
	UserID      int32      `json:"userId"`
	SessionHash string     `json:"sessionHash"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}
