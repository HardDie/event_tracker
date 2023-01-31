package service

import (
	"context"
	"fmt"

	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/repository"
	"github.com/HardDie/event_tracker/internal/utils"
)

type IUser interface {
	Get(ctx context.Context, id, userID int32) (*entity.User, error)

	Password(ctx context.Context, req *dto.UpdatePasswordDTO, userID int32) error
	UpdateProfile(ctx context.Context, req *dto.UpdateProfileDTO) (*entity.User, error)
	UpdateImage(ctx context.Context, req *dto.UpdateProfileImageDTO) (*entity.User, error)
}

type User struct {
	userRepository     repository.IUser
	passwordRepository repository.IPassword

	db *db.DB
}

func NewUser(db *db.DB, repository repository.IUser, password repository.IPassword) *User {
	return &User{
		db:                 db,
		userRepository:     repository,
		passwordRepository: password,
	}
}

func (s *User) Get(ctx context.Context, id, userID int32) (*entity.User, error) {
	return s.userRepository.GetByID(s.db.DB, ctx, id, id == userID)
}

func (s *User) Password(ctx context.Context, req *dto.UpdatePasswordDTO, userID int32) error {
	// Get password from DB
	password, err := s.passwordRepository.GetByUserID(s.db.DB, ctx, userID)
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
	password, err = s.passwordRepository.Update(s.db.DB, ctx, userID, hashPassword)
	if err != nil {
		return err
	}
	return nil
}
func (s *User) UpdateProfile(ctx context.Context, req *dto.UpdateProfileDTO) (*entity.User, error) {
	return s.userRepository.UpdateProfile(s.db.DB, ctx, req)
}
func (s *User) UpdateImage(ctx context.Context, req *dto.UpdateProfileImageDTO) (*entity.User, error) {
	return s.userRepository.UpdateImage(s.db.DB, ctx, req)
}
