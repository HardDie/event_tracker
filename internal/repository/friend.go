package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/HardDie/godb/v2"
	"github.com/dimonrus/gosql"

	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
)

type IFriend interface {
	CreateInvite(tx godb.Queryer, ctx context.Context, userID, id int32) (*entity.FriendInvite, error)
	ListPendingInvitations(tx godb.Queryer, ctx context.Context, userID int32) ([]*dto.InviteListResponseDTO, int32, error)
	DeleteInvite(tx godb.Queryer, ctx context.Context, userID, id int32) error
	GetInviteByUserID(tx godb.Queryer, ctx context.Context, userID, withUserID int32) (*entity.FriendInvite, error)
	GetInviteByID(tx godb.Queryer, ctx context.Context, userID, id int32) (*entity.FriendInvite, error)

	GetFriendByUserID(tx godb.Queryer, ctx context.Context, userID, withUserID int32) (*entity.Friend, error)
	CreateFriendshipLink(tx godb.Queryer, ctx context.Context, userID, withUserID int32) ([]*entity.Friend, error)
	ListOfFriends(tx godb.Queryer, ctx context.Context, userID int32) ([]*entity.User, int32, error)
}

type Friend struct {
}

func NewFriend() *Friend {
	return &Friend{}
}

func (r *Friend) CreateInvite(tx godb.Queryer, ctx context.Context, userID, withUserID int32) (*entity.FriendInvite, error) {
	invite := &entity.FriendInvite{
		UserID:     userID,
		WithUserID: withUserID,
	}

	q := gosql.NewInsert().Into("friend_invites")
	q.Columns().Add("user_id", "with_user_id")
	q.Columns().Arg(userID, withUserID)
	q.Returning().Add("id", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&invite.ID, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return invite, nil
}
func (r *Friend) ListPendingInvitations(tx godb.Queryer, ctx context.Context, userID int32) ([]*dto.InviteListResponseDTO, int32, error) {
	var res []*dto.InviteListResponseDTO

	q := gosql.NewSelect().From("friend_invites fi")
	q.Columns().Add("fi.id", "fi.user_id", "u.displayed_name", "u.profile_image", "fi.created_at")
	q.Relate("JOIN users u ON fi.user_id = u.id")
	q.Where().AddExpression("fi.with_user_id = ?", userID)
	q.Where().AddExpression("fi.deleted_at IS NULL")
	q.Where().AddExpression("u.deleted_at IS NULL")
	q.AddOrder("fi.id")
	rows, err := tx.QueryContext(ctx, q.String(), q.GetArguments()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		friendRequest := &dto.InviteListResponseDTO{}
		err = rows.Scan(&friendRequest.ID, &friendRequest.User.ID, &friendRequest.User.DisplayedName, &friendRequest.User.ProfileImage, &friendRequest.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, friendRequest)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	return res, int32(len(res)), nil
}
func (r *Friend) DeleteInvite(tx godb.Queryer, ctx context.Context, userID, inviteID int32) error {
	q := gosql.NewUpdate().Table("friend_invites")
	q.Set().Add("deleted_at = now()")
	q.Where().AddExpression("id = ?", inviteID)
	q.Where().AddExpression("with_user_id = ? OR user_id = ?", userID, userID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("id")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&inviteID)
	if err != nil {
		return err
	}
	return nil
}
func (r *Friend) GetInviteByUserID(tx godb.Queryer, ctx context.Context, userID, withUserID int32) (*entity.FriendInvite, error) {
	invite := &entity.FriendInvite{
		UserID:     userID,
		WithUserID: withUserID,
	}

	q := gosql.NewSelect().From("friend_invites")
	q.Columns().Add("id", "created_at", "updated_at")
	q.Where().AddExpression("deleted_at IS NULL")
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("with_user_id = ?", withUserID)
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&invite.ID, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return invite, nil
}
func (r *Friend) GetInviteByID(tx godb.Queryer, ctx context.Context, userID, id int32) (*entity.FriendInvite, error) {
	invite := &entity.FriendInvite{
		ID:         id,
		WithUserID: userID,
	}

	q := gosql.NewSelect().From("friend_invites")
	q.Columns().Add("user_id", "created_at", "updated_at")
	q.Where().AddExpression("deleted_at IS NULL")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("with_user_id = ?", userID)
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&invite.UserID, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return invite, nil
}

func (r *Friend) GetFriendByUserID(tx godb.Queryer, ctx context.Context, userID, withUserID int32) (*entity.Friend, error) {
	invite := &entity.Friend{
		UserID:     userID,
		WithUserID: withUserID,
	}

	q := gosql.NewSelect().From("friends")
	q.Columns().Add("id", "created_at", "updated_at")
	q.Where().AddExpression("deleted_at IS NULL")
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("with_user_id = ?", withUserID)
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&invite.ID, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return invite, nil
}
func (r *Friend) CreateFriendshipLink(tx godb.Queryer, ctx context.Context, userID, withUserID int32) ([]*entity.Friend, error) {
	var res []*entity.Friend

	q := gosql.NewInsert().Into("friends")
	q.Columns().Add("user_id", "with_user_id")
	q.Columns().Arg(userID, withUserID)
	q.Columns().Arg(withUserID, userID)
	q.Returning().Add("id", "user_id", "with_user_id", "created_at", "updated_at")
	rows, err := tx.QueryContext(ctx, q.String(), q.GetArguments()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		friend := &entity.Friend{}
		err = rows.Scan(&friend.ID, &friend.UserID, &friend.WithUserID, &friend.CreatedAt, &friend.UpdatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, friend)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (r *Friend) ListOfFriends(tx godb.Queryer, ctx context.Context, userID int32) ([]*entity.User, int32, error) {
	var res []*entity.User

	q := gosql.NewSelect().From("friends f")
	q.Columns().Add("u.id", "u.displayed_name", "u.profile_image")
	q.Relate("JOIN users u ON f.with_user_id = u.id")
	q.Where().AddExpression("f.user_id = ?", userID)
	q.Where().AddExpression("f.deleted_at IS NULL")
	q.Where().AddExpression("u.deleted_at IS NULL")
	q.AddOrder("f.id")
	rows, err := tx.QueryContext(ctx, q.String(), q.GetArguments()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &entity.User{}
		err = rows.Scan(&user.ID, &user.DisplayedName, &user.ProfileImage)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	return res, int32(len(res)), nil
}
