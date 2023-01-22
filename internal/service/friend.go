package service

import (
	"context"
	"errors"

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
}

func NewFriend(repository repository.IFriend, user repository.IUser) *Friend {
	return &Friend{
		repository:     repository,
		userRepository: user,
	}
}

func (s *Friend) InviteFriend(ctx context.Context, req *dto.InviteFriendDTO) error {
	// Check if such a user exists
	user, err := s.userRepository.GetByName(ctx, req.Username)
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
	friend, err := s.repository.GetFriendByUserID(ctx, req.ID, user.ID)
	if err != nil {
		return err
	}
	if friend != nil {
		return errors.New("already friends")
	}

	// Check if such an invitation already exists
	inv, err := s.repository.GetInviteByUserID(ctx, req.ID, user.ID)
	if err != nil {
		return err
	}
	if inv != nil {
		return errors.New("such invitation already exist")
	}

	// Create an invitation for a friend
	_, err = s.repository.CreateInvite(ctx, req.ID, user.ID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Friend) InviteListPending(ctx context.Context, userID int32) ([]*dto.InviteListResponseDTO, int32, error) {
	return s.repository.ListPendingInvitations(ctx, userID)
}
func (s *Friend) AcceptFriendship(ctx context.Context, userID, inviteID int32) error {
	// Get information about the invitation, check if there is such an invitation
	invite, err := s.repository.GetInviteByID(ctx, userID, inviteID)
	if err != nil {
		return err
	}

	// Check if the invitation is associated with the current user
	if invite.WithUserID != userID {
		return errors.New("bad invite")
	}

	// Delete invitation
	err = s.repository.DeleteInvite(ctx, userID, invite.ID)
	if err != nil {
		return err
	}

	// Create a friendship
	_, err = s.repository.CreateFriendshipLink(ctx, userID, invite.UserID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Friend) RejectFriendship(ctx context.Context, userID, inviteID int32) error {
	// Get information about the invitation, check if there is such an invitation
	invite, err := s.repository.GetInviteByID(ctx, userID, inviteID)
	if err != nil {
		return err
	}

	// Check if the invitation is associated with the current user
	if invite.WithUserID != userID {
		return errors.New("bad invite")
	}

	// Delete invitation
	err = s.repository.DeleteInvite(ctx, userID, invite.ID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Friend) ListOfFriends(ctx context.Context, userID int32) ([]*entity.User, int32, error) {
	return s.repository.ListOfFriends(ctx, userID)
}
