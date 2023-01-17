package entity

import "time"

type EventType struct {
	ID        int32      `json:"id"`
	UserID    int32      `json:"userId"`
	EventType string     `json:"eventType"`
	IsVisible bool       `json:"isVisible"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
