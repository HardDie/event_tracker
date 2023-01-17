package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
)

type IUser interface {
	GetByID(ctx context.Context, id int32, showPrivateInfo bool) (*entity.User, error)
	GetByName(ctx context.Context, name string) (*entity.User, error)
	Create(ctx context.Context, name, displayedName string) (*entity.User, error)
	Update(ctx context.Context, req *dto.UpdateProfileDTO, id int32) (*entity.User, error)
}

type User struct {
	db *db.DB
}

func NewUser(db *db.DB) *User {
	return &User{
		db: db,
	}
}

func (r *User) GetByID(ctx context.Context, id int32, showPrivateInfo bool) (*entity.User, error) {
	user := &entity.User{
		ID: id,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("displayed_name", "created_at", "updated_at", "deleted_at")
	if showPrivateInfo {
		q.Columns().Add("username", "email")
	}
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	var err error
	if showPrivateInfo {
		err = row.Scan(&user.DisplayedName, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
			&user.Username, &user.Email)
	} else {
		err = row.Scan(&user.DisplayedName, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil

}
func (r *User) GetByName(ctx context.Context, name string) (*entity.User, error) {
	user := &entity.User{
		Username: name,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("id", "displayed_name", "email", "created_at", "updated_at", "deleted_at")
	q.Where().AddExpression("username = ?", name)
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&user.ID, &user.DisplayedName, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (r *User) Create(ctx context.Context, name, displayedName string) (*entity.User, error) {
	user := &entity.User{
		Username:      name,
		DisplayedName: displayedName,
	}

	q := gosql.NewInsert().Into("users")
	q.Columns().Add("username", "displayed_name")
	q.Columns().Arg(name, displayedName)
	q.Returning().Add("id", "created_at", "updated_at")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *User) Update(ctx context.Context, req *dto.UpdateProfileDTO, id int32) (*entity.User, error) {
	user := &entity.User{
		ID:            id,
		DisplayedName: req.DisplayedName,
		Email:         req.Email,
	}

	q := gosql.NewUpdate().Table("users")
	q.Set().Append("displayed_name = ?", req.DisplayedName)
	q.Set().Append("email = ?", req.Email)
	q.Set().Append("updated_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("username", "created_at", "updated_at")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
