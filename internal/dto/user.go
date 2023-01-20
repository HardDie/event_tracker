package dto

type GetUserDTO struct {
	ID int32 `json:"id" validate:"gt=0"`
}

type UpdatePasswordDTO struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,nefield=OldPassword"`
}

type UpdateProfileDTO struct {
	ID            int32   `json:"-" validate:"gt=0"`
	DisplayedName string  `json:"displayedName" validate:"required"`
	Email         *string `json:"email" validate:"omitempty,email"`
}

type UpdateProfileImageDTO struct {
	ID           int32   `json:"-" validate:"gt=0"`
	ProfileImage *string `json:"profileImage" validate:"omitempty,max=10000,base64"`
}
