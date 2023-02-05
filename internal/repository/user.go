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

type IUser interface {
	GetByID(tx godb.Queryer, ctx context.Context, id int32, showPrivateInfo bool) (*entity.User, error)
	GetByName(tx godb.Queryer, ctx context.Context, name string) (*entity.User, error)
	Create(tx godb.Queryer, ctx context.Context, name, displayedName string) (*entity.User, error)
	UpdateProfile(tx godb.Queryer, ctx context.Context, req *dto.UpdateProfileDTO) (*entity.User, error)
	UpdateImage(tx godb.Queryer, ctx context.Context, req *dto.UpdateProfileImageDTO) (*entity.User, error)
}

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (r *User) GetByID(tx godb.Queryer, ctx context.Context, id int32, showPrivateInfo bool) (*entity.User, error) {
	user := &entity.User{
		ID: id,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("displayed_name", "profile_image", "created_at", "updated_at", "deleted_at")
	if showPrivateInfo {
		q.Columns().Add("username", "email")
	}
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	var err error
	if showPrivateInfo {
		err = row.Scan(&user.DisplayedName, &user.ProfileImage, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
			&user.Username, &user.Email)
	} else {
		err = row.Scan(&user.DisplayedName, &user.ProfileImage, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil

}
func (r *User) GetByName(tx godb.Queryer, ctx context.Context, name string) (*entity.User, error) {
	user := &entity.User{
		Username: name,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("id", "displayed_name", "email", "profile_image", "created_at", "updated_at", "deleted_at")
	q.Where().AddExpression("username = ?", name)
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&user.ID, &user.DisplayedName, &user.Email, &user.ProfileImage, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (r *User) Create(tx godb.Queryer, ctx context.Context, name, displayedName string) (*entity.User, error) {
	user := &entity.User{
		Username:      name,
		DisplayedName: displayedName,
	}

	q := gosql.NewInsert().Into("users")
	q.Columns().Add("username", "displayed_name")
	q.Columns().Arg(name, displayedName)
	q.Returning().Add("id", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *User) UpdateProfile(tx godb.Queryer, ctx context.Context, req *dto.UpdateProfileDTO) (*entity.User, error) {
	user := &entity.User{
		ID:            req.ID,
		DisplayedName: req.DisplayedName,
		Email:         req.Email,
	}

	q := gosql.NewUpdate().Table("users")
	q.Set().Append("displayed_name = ?", req.DisplayedName)
	q.Set().Append("email = ?", req.Email)
	q.Set().Append("updated_at = now()")
	q.Where().AddExpression("id = ?", req.ID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("username", "profile_image", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&user.Username, &user.ProfileImage, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *User) UpdateImage(tx godb.Queryer, ctx context.Context, req *dto.UpdateProfileImageDTO) (*entity.User, error) {
	user := &entity.User{
		ID:           req.ID,
		ProfileImage: req.ProfileImage,
	}

	q := gosql.NewUpdate().Table("users")
	q.Set().Append("profile_image = ?", req.ProfileImage)
	q.Set().Append("updated_at = now()")
	q.Where().AddExpression("id = ?", req.ID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("username", "displayed_name", "email", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&user.Username, &user.DisplayedName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
