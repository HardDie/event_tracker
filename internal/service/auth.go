package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HardDie/event_tracker/internal/config"
	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/logger"
	"github.com/HardDie/event_tracker/internal/repository"
	"github.com/HardDie/event_tracker/internal/utils"
)

var (
	ErrorSessionHasExpired = errors.New("session has expired")
)

type IAuth interface {
	Register(ctx context.Context, req *dto.RegisterDTO) (*entity.User, error)
	Login(ctx context.Context, req *dto.LoginDTO) (*entity.User, error)
	Logout(ctx context.Context, sessionID int32) error
	GenerateCookie(ctx context.Context, userID int32) (*entity.Session, error)
	ValidateCookie(ctx context.Context, session string) (*entity.Session, error)
	GetUserInfo(ctx context.Context, userID int32) (*entity.User, error)
}

type Auth struct {
	userRepository     repository.IUser
	passwordRepository repository.IPassword
	sessionRepository  repository.ISession

	cfg *config.Config
	db  *db.DB
}

func NewAuth(db *db.DB, cfg *config.Config, user repository.IUser, password repository.IPassword,
	session repository.ISession) *Auth {
	return &Auth{
		db:                 db,
		cfg:                cfg,
		userRepository:     user,
		passwordRepository: password,
		sessionRepository:  session,
	}
}

func (s *Auth) Register(ctx context.Context, req *dto.RegisterDTO) (*entity.User, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Check if username is not busy
	user, err := s.userRepository.GetByName(tx, ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("error while trying get user: %w", err)
	}
	if user != nil {
		return nil, fmt.Errorf("username already exist")
	}

	// Hashing password
	hashPassword, err := utils.HashBcrypt(req.Password)
	if err != nil {
		return nil, err
	}

	// Create a user
	user, err = s.userRepository.Create(tx, ctx, req.Username, req.DisplayedName)
	if err != nil {
		return nil, err
	}

	// Create a password
	_, err = s.passwordRepository.Create(tx, ctx, user.ID, hashPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *Auth) Login(ctx context.Context, req *dto.LoginDTO) (*entity.User, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Check if such user exist
	user, err := s.userRepository.GetByName(tx, ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not exist")
	}

	// Get password from DB
	password, err := s.passwordRepository.GetByUserID(tx, ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if password == nil {
		return nil, fmt.Errorf("password for user not exist")
	}

	// Check if the password is locked after failed attempts
	if password.FailedAttempts >= int32(s.cfg.PwdMaxAttempts) {
		// Check if the password block time has expired
		if time.Now().Sub(password.UpdatedAt) <= time.Hour*time.Duration(s.cfg.PwdBlockTime) {
			return nil, fmt.Errorf("user was blocked after failed attempts")
		}
		// If the blocking time has expired, reset the counter of failed attempts
		password, err = s.passwordRepository.ResetFailedAttempts(tx, ctx, password.ID)
		if err != nil {
			return nil, fmt.Errorf("error resetting the counter of failed attempts: %w", err)
		}
	}

	// Check if password is correct
	if !utils.HashBcryptCompare(req.Password, password.PasswordHash) {
		// Increased number of failed attempts
		_, err = s.passwordRepository.IncreaseFailedAttempts(tx, ctx, password.ID)
		if err != nil {
			logger.Error.Println("Error increasing failed attempts:", err.Error())
		}
		return nil, fmt.Errorf("invalid password")
	}

	// Reset the failed attempts counter after the first successful attempt
	if password.FailedAttempts > 0 {
		_, err = s.passwordRepository.ResetFailedAttempts(tx, ctx, password.ID)
		if err != nil {
			logger.Error.Println("Error flushing failed attempts:", err.Error())
		}
	}
	return user, nil
}
func (s *Auth) Logout(ctx context.Context, sessionID int32) error {
	return s.sessionRepository.DeleteByID(s.db.DB, ctx, sessionID)
}
func (s *Auth) GenerateCookie(ctx context.Context, userID int32) (*entity.Session, error) {
	// Generate session key
	sessionKey, err := utils.GenerateSessionKey()
	if err != nil {
		return nil, fmt.Errorf("generate session key: %w", err)
	}

	// Write session to DB
	resp, err := s.sessionRepository.CreateOrUpdate(s.db.DB, ctx, userID, utils.HashSha256(sessionKey))
	if err != nil {
		return nil, fmt.Errorf("write session to DB: %w", err)
	}
	resp.Session = sessionKey

	return resp, nil
}
func (s *Auth) ValidateCookie(ctx context.Context, sessionToken string) (*entity.Session, error) {
	// Check if session exist
	sessionHash := utils.HashSha256(sessionToken)
	session, err := s.sessionRepository.GetByUserID(s.db.DB, ctx, sessionHash)
	if err != nil {
		return nil, err
	}

	// Check if session is not expired
	if time.Now().Sub(session.UpdatedAt) > time.Hour*24 {
		return nil, ErrorSessionHasExpired
	}
	return session, nil
}
func (s *Auth) GetUserInfo(ctx context.Context, userID int32) (*entity.User, error) {
	return s.userRepository.GetByID(s.db.DB, ctx, userID, true)
}
