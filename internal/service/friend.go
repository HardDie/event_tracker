package service

import (
	"context"

	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/errs"
	"github.com/HardDie/event_tracker/internal/logger"
	"github.com/HardDie/event_tracker/internal/repository"
)

type IFriend interface {
	InviteFriend(ctx context.Context, req *dto.InviteFriendDTO) error
	InviteListPending(ctx context.Context, userID int32) ([]*dto.InviteListResponseDTO, int32, error)
	AcceptFriendship(ctx context.Context, userID, inviteID int32) error
	RejectFriendship(ctx context.Context, userID, inviteID int32) error
	ListOfFriends(ctx context.Context, userID int32) ([]*entity.User, int32, error)
}

type Friend struct {
	repository     repository.IFriend
	userRepository repository.IUser

	db *db.DB
}

func NewFriend(db *db.DB, repository repository.IFriend, user repository.IUser) *Friend {
	return &Friend{
		db:             db,
		repository:     repository,
		userRepository: user,
	}
}

func (s *Friend) InviteFriend(ctx context.Context, req *dto.InviteFriendDTO) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		logger.Error.Printf("error starting transaction: %v", err.Error())
		return errs.InternalError
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Check if such a user exists
	user, err := s.userRepository.GetByName(tx, ctx, req.Username)
	if err != nil {
		logger.Error.Printf("error while trying get user: %v", err.Error())
		return errs.InternalError
	}
	if user == nil {
		return errs.BadRequest.AddMessage("invalid username, such user not exist")
	}

	// Check if the user is trying to add himself to the friends list
	if user.ID == req.ID {
		return errs.BadRequest.AddMessage("can't invite your own account")
	}

	// Check if the user is already a friend
	friend, err := s.repository.GetFriendByUserID(tx, ctx, req.ID, user.ID)
	if err != nil {
		logger.Error.Printf("error trying get friend: %v", err.Error())
		return errs.InternalError
	}
	if friend != nil {
		return errs.BadRequest.AddMessage("already friends")
	}

	// Check if such an invitation already exists
	inv, err := s.repository.GetInviteByUserID(tx, ctx, req.ID, user.ID)
	if err != nil {
		logger.Error.Printf("error trying get invitation: %v", err.Error())
		return errs.InternalError
	}
	if inv != nil {
		return errs.BadRequest.AddMessage("such invitation already exist")
	}

	// Create an invitation for a friend
	_, err = s.repository.CreateInvite(tx, ctx, req.ID, user.ID)
	if err != nil {
		logger.Error.Printf("error while creating invitation: %v", err.Error())
		return errs.InternalError
	}
	return nil
}
func (s *Friend) InviteListPending(ctx context.Context, userID int32) ([]*dto.InviteListResponseDTO, int32, error) {
	res, cnt, err := s.repository.ListPendingInvitations(s.db.DB, ctx, userID)
	if err != nil {
		logger.Error.Printf("error list pending invites: %v", err.Error())
		return nil, 0, errs.InternalError
	}
	return res, cnt, nil
}
func (s *Friend) AcceptFriendship(ctx context.Context, userID, inviteID int32) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		logger.Error.Printf("error starting transaction: %v", err.Error())
		return errs.InternalError
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Get information about the invitation, check if there is such an invitation
	invite, err := s.repository.GetInviteByID(tx, ctx, userID, inviteID)
	if err != nil {
		logger.Error.Printf("error trying get invitation: %v", err.Error())
		return errs.InternalError
	}
	if invite == nil {
		return errs.BadRequest.AddMessage("there is no such invitation")
	}

	// Check if the invitation is associated with the current user
	if invite.WithUserID != userID {
		logger.Error.Printf("invalid invite: %v", invite.ID)
		return errs.InternalError
	}

	// Delete invitation
	err = s.repository.DeleteInvite(tx, ctx, userID, invite.ID)
	if err != nil {
		logger.Error.Printf("error trying delete invitation: %v", err.Error())
		return errs.InternalError
	}

	// Check if there exist is an outgoing friendship invitation
	outgoingInvite, err := s.repository.GetInviteByUserID(tx, ctx, userID, invite.WithUserID)
	if err != nil {
		logger.Error.Printf("error trying get invitation: %v", err.Error())
		return errs.InternalError
	}
	if outgoingInvite != nil {
		// If an invitation exists, delete it
		err = s.repository.DeleteInvite(tx, ctx, userID, outgoingInvite.ID)
		if err != nil {
			logger.Error.Printf("error trying delete invitation: %v", err.Error())
			return errs.InternalError
		}
	}

	// Create a friendship
	_, err = s.repository.CreateFriendshipLink(tx, ctx, userID, invite.UserID)
	if err != nil {
		logger.Error.Printf("error while creating friend link: %v", err.Error())
		return errs.InternalError
	}
	return nil
}
func (s *Friend) RejectFriendship(ctx context.Context, userID, inviteID int32) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		logger.Error.Printf("error starting transaction: %v", err.Error())
		return errs.InternalError
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Get information about the invitation, check if there is such an invitation
	invite, err := s.repository.GetInviteByID(tx, ctx, userID, inviteID)
	if err != nil {
		logger.Error.Printf("error trying get invitation: %v", err.Error())
		return errs.InternalError
	}
	if invite == nil {
		return errs.BadRequest.AddMessage("there is no such invitation")
	}

	// Check if the invitation is associated with the current user
	if invite.WithUserID != userID {
		logger.Error.Printf("invalid invite: %v", invite.ID)
		return errs.InternalError
	}

	// Delete invitation
	err = s.repository.DeleteInvite(tx, ctx, userID, invite.ID)
	if err != nil {
		logger.Error.Printf("error trying delete invitation: %v", err.Error())
		return errs.InternalError
	}
	return nil
}
func (s *Friend) ListOfFriends(ctx context.Context, userID int32) ([]*entity.User, int32, error) {
	res, cnt, err := s.repository.ListOfFriends(s.db.DB, ctx, userID)
	if err != nil {
		logger.Error.Printf("error list of friends: %v", err.Error())
		return nil, 0, errs.InternalError
	}
	return res, cnt, nil
}
