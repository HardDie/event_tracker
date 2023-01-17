package dto

type RegisterDTO struct {
	Username      string `json:"username" validate:"required"`
	Password      string `json:"password" validate:"required"`
	DisplayedName string `json:"displayedName" validate:"required"`
}

type LoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
