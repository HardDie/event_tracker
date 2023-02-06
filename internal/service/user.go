package service

import (
	"context"

	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/errs"
	"github.com/HardDie/event_tracker/internal/logger"
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
	user, err := s.userRepository.GetByID(s.db.DB, ctx, id, id == userID)
	if err != nil {
		logger.Error.Printf("error get user: %v", err.Error())
		return nil, errs.InternalError
	}
	return user, nil
}

func (s *User) Password(ctx context.Context, req *dto.UpdatePasswordDTO, userID int32) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		logger.Error.Printf("error starting transaction: %v", err.Error())
		return errs.InternalError
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Get password from DB
	password, err := s.passwordRepository.GetByUserID(tx, ctx, userID)
	if err != nil {
		logger.Error.Printf("error read password from DB: %v", err.Error())
		return errs.InternalError
	}
	if password == nil {
		logger.Error.Printf("password for user %d not found in DB", userID)
		return errs.InternalError
	}

	// Check if password is correct
	if !utils.HashBcryptCompare(req.OldPassword, password.PasswordHash) {
		return errs.BadRequest.AddMessage("invalid old password")
	}

	// Hashing password
	hashPassword, err := utils.HashBcrypt(req.NewPassword)
	if err != nil {
		logger.Error.Printf("error hashing password: %v", err.Error())
		return errs.InternalError
	}

	// Update password
	password, err = s.passwordRepository.Update(tx, ctx, userID, hashPassword)
	if err != nil {
		logger.Error.Printf("error updating password in DB: %v", err.Error())
		return errs.InternalError
	}
	return nil
}
func (s *User) UpdateProfile(ctx context.Context, req *dto.UpdateProfileDTO) (*entity.User, error) {
	user, err := s.userRepository.UpdateProfile(s.db.DB, ctx, req)
	if err != nil {
		logger.Error.Printf("error update user profile: %v", err.Error())
		return nil, errs.InternalError
	}
	return user, nil
}
func (s *User) UpdateImage(ctx context.Context, req *dto.UpdateProfileImageDTO) (*entity.User, error) {
	user, err := s.userRepository.UpdateImage(s.db.DB, ctx, req)
	if err != nil {
		logger.Error.Printf("error update user image: %v", err.Error())
		return nil, errs.InternalError
	}
	return user, nil
}
