package service

import (
	"context"
	"fmt"

	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/repository"
	"github.com/HardDie/event_tracker/internal/utils"
)

type IUser interface {
	Get(ctx context.Context, id int32) (*entity.User, error)

	Password(ctx context.Context, req *dto.UpdatePasswordDTO, userID int32) error
	Profile(ctx context.Context, req *dto.UpdateProfileDTO, userID int32) (*entity.User, error)
}

type User struct {
	userRepository     repository.IUser
	passwordRepository repository.IPassword
}

func NewUser(repository repository.IUser, password repository.IPassword) *User {
	return &User{
		userRepository:     repository,
		passwordRepository: password,
	}
}

func (s *User) Get(ctx context.Context, id int32) (*entity.User, error) {
	return s.userRepository.GetByID(ctx, id, false)
}

func (s *User) Password(ctx context.Context, req *dto.UpdatePasswordDTO, userID int32) error {
	// Get password from DB
	password, err := s.passwordRepository.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if password == nil {
		return fmt.Errorf("password for user not exist")
	}

	// Check if password is correct
	if !utils.HashBcryptCompare(req.OldPassword, password.PasswordHash) {
		return fmt.Errorf("invalid old password")
	}

	// Hashing password
	hashPassword, err := utils.HashBcrypt(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	password, err = s.passwordRepository.Update(ctx, userID, hashPassword)
	if err != nil {
		return err
	}
	return nil
}
func (s *User) Profile(ctx context.Context, req *dto.UpdateProfileDTO, userID int32) (*entity.User, error) {
	return s.userRepository.Update(ctx, req, userID)
}
