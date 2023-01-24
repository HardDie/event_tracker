package dto

import "time"

type CreateEventTypeDTO struct {
	Name      string `json:"name" validate:"required"`
	IsVisible bool   `json:"isVisible"`
}

type EditEventTypeDTO struct {
	ID        int32  `json:"-" validate:"required,gt=0"`
	Name      string `json:"name" validate:"required"`
	IsVisible bool   `json:"isVisible"`
}

type CreateEventDTO struct {
	EventTypeID int32     `json:"eventTypeId" validate:"required,gt=0"`
	Date        time.Time `json:"date" validate:"required"`
}
type ListEventDTO struct {
	UserID     *int32     `json:"userId" validate:"omitempty,gt=0"`
	TypeID     *int32     `json:"typeId" validate:"omitempty,gt=0"`
	PeriodType PeriodType `json:"periodType" validate:"required,gt=0,lt=4"`
	Date       time.Time  `json:"date" validate:"required"`
}

type FeedResponseDTO struct {
	UserID      int32     `json:"userId"`
	EventTypeID int32     `json:"eventTypeId"`
	EventType   string    `json:"eventType"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"createdAt"`
}

/*
 * internal
 */

type PeriodType int32

const (
	PeriodDay   PeriodType = 1
	PeriodMonth            = 2
	PeriodYear             = 3
)

type ListEventFilter struct {
	UserID      int32
	TypeID      *int32
	OnlyVisible bool
	PeriodType  PeriodType
	Date        time.Time
}
