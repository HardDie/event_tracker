package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
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
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Check if such a user exists
	user, err := s.userRepository.GetByName(tx, ctx, req.Username)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid username, such user not exist")
	}

	// Check if the user is trying to add himself to the friends list
	if user.ID == req.ID {
		return errors.New("can't invite your own account")
	}

	// Check if the user is already a friend
	friend, err := s.repository.GetFriendByUserID(tx, ctx, req.ID, user.ID)
	if err != nil {
		return err
	}
	if friend != nil {
		return errors.New("already friends")
	}

	// Check if such an invitation already exists
	inv, err := s.repository.GetInviteByUserID(tx, ctx, req.ID, user.ID)
	if err != nil {
		return err
	}
	if inv != nil {
		return errors.New("such invitation already exist")
	}

	// Create an invitation for a friend
	_, err = s.repository.CreateInvite(tx, ctx, req.ID, user.ID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Friend) InviteListPending(ctx context.Context, userID int32) ([]*dto.InviteListResponseDTO, int32, error) {
	return s.repository.ListPendingInvitations(s.db.DB, ctx, userID)
}
func (s *Friend) AcceptFriendship(ctx context.Context, userID, inviteID int32) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Get information about the invitation, check if there is such an invitation
	invite, err := s.repository.GetInviteByID(tx, ctx, userID, inviteID)
	if err != nil {
		return err
	}
	if invite == nil {
		return errors.New("there is no such invitation")
	}

	// Check if the invitation is associated with the current user
	if invite.WithUserID != userID {
		return errors.New("bad invite")
	}

	// Delete invitation
	err = s.repository.DeleteInvite(tx, ctx, userID, invite.ID)
	if err != nil {
		return err
	}

	// Check if there exist is an outgoing friendship invitation
	outgoingInvite, err := s.repository.GetInviteByUserID(tx, ctx, userID, invite.WithUserID)
	if err != nil {
		return err
	}
	if outgoingInvite != nil {
		// If an invitation exists, delete it
		err = s.repository.DeleteInvite(tx, ctx, userID, outgoingInvite.ID)
		if err != nil {
			return err
		}
	}

	// Create a friendship
	_, err = s.repository.CreateFriendshipLink(tx, ctx, userID, invite.UserID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Friend) RejectFriendship(ctx context.Context, userID, inviteID int32) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { s.db.EndTx(tx, err) }()

	// Get information about the invitation, check if there is such an invitation
	invite, err := s.repository.GetInviteByID(tx, ctx, userID, inviteID)
	if err != nil {
		return err
	}
	if invite == nil {
		return errors.New("there is no such invitation")
	}

	// Check if the invitation is associated with the current user
	if invite.WithUserID != userID {
		return errors.New("bad invite")
	}

	// Delete invitation
	err = s.repository.DeleteInvite(tx, ctx, userID, invite.ID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Friend) ListOfFriends(ctx context.Context, userID int32) ([]*entity.User, int32, error) {
	return s.repository.ListOfFriends(s.db.DB, ctx, userID)
}
