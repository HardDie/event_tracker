package dto

import (
	"time"

	"github.com/HardDie/event_tracker/internal/entity"
)

type InviteFriendDTO struct {
	ID       int32  `json:"-" validate:"gt=0"`
	Username string `json:"username" validate:"required"`
}

type InviteListResponseDTO struct {
	ID        int32       `json:"id"`
	User      entity.User `json:"user"`
	CreatedAt time.Time   `json:"createdAt"`
}
