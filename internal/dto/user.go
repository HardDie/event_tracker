package dto

type GetUserDTO struct {
	ID int32 `json:"id" validate:"gt=0"`
}

type UpdatePasswordDTO struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,nefield=OldPassword"`
}

type UpdateProfileDTO struct {
	DisplayedName string  `json:"displayedName" validate:"required"`
	Email         *string `json:"email" validate:"omitempty,email"`
}
