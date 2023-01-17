package entity

import "time"

type Event struct {
	ID        int32      `json:"id"`
	UserID    int32      `json:"userId"`
	TypeID    int32      `json:"eventTypeId"`
	Date      time.Time  `json:"date"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
