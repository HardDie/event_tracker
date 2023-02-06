package service

import (
	"context"
	"time"

	"github.com/HardDie/event_tracker/internal/config"
	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/errs"
	"github.com/HardDie/event_tracker/internal/logger"
	"github.com/HardDie/event_tracker/internal/repository"
	"github.com/HardDie/event_tracker/internal/utils"
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
		logger.Error.Printf("error starting transaction: %v", err.Error())
		return nil, errs.InternalError
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Check if username is not busy
	user, err := s.userRepository.GetByName(tx, ctx, req.Username)
	if err != nil {
		logger.Error.Printf("error while trying get user: %v", err.Error())
		return nil, errs.InternalError
	}
	if user != nil {
		return nil, errs.BadRequest.AddMessage("username already exist")
	}

	// Hashing password
	hashPassword, err := utils.HashBcrypt(req.Password)
	if err != nil {
		logger.Error.Printf("error hash bcrypt: %v", err.Error())
		return nil, errs.InternalError
	}

	// Create a user
	user, err = s.userRepository.Create(tx, ctx, req.Username, req.DisplayedName)
	if err != nil {
		logger.Error.Printf("error writing user into DB: %v", err.Error())
		return nil, errs.InternalError
	}

	// Create a password
	_, err = s.passwordRepository.Create(tx, ctx, user.ID, hashPassword)
	if err != nil {
		logger.Error.Printf("error writing password into DB: %v", err.Error())
		return nil, errs.InternalError
	}

	return user, nil
}
func (s *Auth) Login(ctx context.Context, req *dto.LoginDTO) (*entity.User, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		logger.Error.Printf("error starting transaction: %v", err.Error())
		return nil, errs.InternalError
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Check if such user exist
	user, err := s.userRepository.GetByName(tx, ctx, req.Username)
	if err != nil {
		logger.Error.Printf("error while trying get user: %v", err.Error())
		return nil, errs.InternalError
	}
	if user == nil {
		return nil, errs.BadRequest.AddMessage("username or password is invalid")
	}

	// Get password from DB
	password, err := s.passwordRepository.GetByUserID(tx, ctx, user.ID)
	if err != nil {
		logger.Error.Printf("error while trying get password: %v", err.Error())
		return nil, errs.InternalError
	}
	if password == nil {
		logger.Error.Printf("password for user %d not found", user.ID)
		return nil, errs.InternalError
	}

	// Check if the password is locked after failed attempts
	if password.FailedAttempts >= int32(s.cfg.PwdMaxAttempts) {
		// Check if the password block time has expired
		if time.Now().Sub(password.UpdatedAt) <= time.Hour*time.Duration(s.cfg.PwdBlockTime) {
			return nil, errs.UserBlocked.AddMessage("too many invalid requests")
		}
		// If the blocking time has expired, reset the counter of failed attempts
		password, err = s.passwordRepository.ResetFailedAttempts(tx, ctx, password.ID)
		if err != nil {
			logger.Error.Printf("error resetting the counter of failed attempts: %v", err)
			return nil, errs.InternalError
		}
	}

	// Check if password is correct
	if !utils.HashBcryptCompare(req.Password, password.PasswordHash) {
		// Increased number of failed attempts
		_, err = s.passwordRepository.IncreaseFailedAttempts(tx, ctx, password.ID)
		if err != nil {
			logger.Error.Printf("Error increasing failed attempts: %v", err.Error())
		}
		return nil, errs.BadRequest.AddMessage("username or password is invalid")
	}

	// Reset the failed attempts counter after the first successful attempt
	if password.FailedAttempts > 0 {
		_, err = s.passwordRepository.ResetFailedAttempts(tx, ctx, password.ID)
		if err != nil {
			logger.Error.Printf("Error flushing failed attempts: %v", err.Error())
		}
	}
	return user, nil
}
func (s *Auth) Logout(ctx context.Context, sessionID int32) error {
	err := s.sessionRepository.DeleteByID(s.db.DB, ctx, sessionID)
	if err != nil {
		logger.Error.Printf("error deleting session: %v", err.Error())
		return errs.InternalError
	}
	return nil
}
func (s *Auth) GenerateCookie(ctx context.Context, userID int32) (*entity.Session, error) {
	// Generate session key
	sessionKey, err := utils.GenerateSessionKey()
	if err != nil {
		logger.Error.Printf("error generate session key: %v", err)
		return nil, errs.InternalError
	}

	// Write session to DB
	resp, err := s.sessionRepository.CreateOrUpdate(s.db.DB, ctx, userID, utils.HashSha256(sessionKey))
	if err != nil {
		logger.Error.Printf("write session to DB: %v", err)
		return nil, errs.InternalError
	}
	resp.Session = sessionKey

	return resp, nil
}
func (s *Auth) ValidateCookie(ctx context.Context, sessionToken string) (*entity.Session, error) {
	// Check if session exist
	sessionHash := utils.HashSha256(sessionToken)
	session, err := s.sessionRepository.GetByUserID(s.db.DB, ctx, sessionHash)
	if err != nil {
		logger.Error.Printf("error read session from db: %v", err.Error())
		return nil, errs.InternalError
	}
	if session == nil {
		return nil, errs.SessionInvalid.AddMessage("session not exist")
	}

	// Check if session is not expired
	if time.Now().Sub(session.UpdatedAt) > time.Hour*24 {
		return nil, errs.SessionInvalid.AddMessage("session has expired")
	}
	return session, nil
}
func (s *Auth) GetUserInfo(ctx context.Context, userID int32) (*entity.User, error) {
	user, err := s.userRepository.GetByID(s.db.DB, ctx, userID, true)
	if err != nil {
		logger.Error.Printf("error get user info: %v", err.Error())
		return nil, errs.InternalError
	}
	return user, nil
}
