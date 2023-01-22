package dto

import "time"

type InviteFriendDTO struct {
	ID       int32  `json:"-" validate:"gt=0"`
	Username string `json:"username" validate:"required"`
}

type InviteListResponseDTO struct {
	ID            int32     `json:"id"`
	UserID        int32     `json:"userId"`
	DisplayedName string    `json:"name"`
	ProfileImage  *string   `json:"profileImage"`
	CreatedAt     time.Time `json:"createdAt"`
}
