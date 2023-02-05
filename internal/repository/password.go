package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/HardDie/godb/v2"
	"github.com/dimonrus/gosql"

	"github.com/HardDie/event_tracker/internal/entity"
)

type IPassword interface {
	Create(tx godb.Queryer, ctx context.Context, userID int32, passwordHash string) (*entity.Password, error)
	GetByUserID(tx godb.Queryer, ctx context.Context, userID int32) (*entity.Password, error)
	Update(tx godb.Queryer, ctx context.Context, id int32, passwordHash string) (*entity.Password, error)
	IncreaseFailedAttempts(tx godb.Queryer, ctx context.Context, id int32) (*entity.Password, error)
	ResetFailedAttempts(tx godb.Queryer, ctx context.Context, id int32) (*entity.Password, error)
}

type Password struct {
}

func NewPassword() *Password {
	return &Password{}
}

func (r *Password) Create(tx godb.Queryer, ctx context.Context, userID int32, passwordHash string) (*entity.Password, error) {
	password := &entity.Password{
		UserID:       userID,
		PasswordHash: passwordHash,
	}

	q := gosql.NewInsert().Into("passwords")
	q.Columns().Add("user_id", "password_hash")
	q.Columns().Arg(userID, passwordHash)
	q.Returning().Add("id", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&password.ID, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return password, nil
}
func (r *Password) GetByUserID(tx godb.Queryer, ctx context.Context, userID int32) (*entity.Password, error) {
	password := &entity.Password{
		UserID: userID,
	}

	q := gosql.NewSelect().From("passwords")
	q.Columns().Add("id", "password_hash", "failed_attempts", "created_at", "updated_at")
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&password.ID, &password.PasswordHash, &password.FailedAttempts, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return password, nil
}
func (r *Password) Update(tx godb.Queryer, ctx context.Context, id int32, passwordHash string) (*entity.Password, error) {
	password := &entity.Password{
		ID:           id,
		PasswordHash: passwordHash,
	}

	q := gosql.NewUpdate().Table("passwords")
	q.Set().Append("password_hash = ?", passwordHash)
	q.Set().Append("updated_at = now()")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("user_id", "failed_attempts", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&password.UserID, &password.FailedAttempts, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return password, nil
}
func (r *Password) IncreaseFailedAttempts(tx godb.Queryer, ctx context.Context, id int32) (*entity.Password, error) {
	password := &entity.Password{
		ID: id,
	}

	q := gosql.NewUpdate().Table("passwords")
	q.Set().Add("failed_attempts = failed_attempts + 1")
	q.Set().Append("updated_at = now()")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("user_id", "password_hash", "failed_attempts", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&password.UserID, &password.PasswordHash, &password.FailedAttempts, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return password, nil
}
func (r *Password) ResetFailedAttempts(tx godb.Queryer, ctx context.Context, id int32) (*entity.Password, error) {
	password := &entity.Password{
		ID: id,
	}

	q := gosql.NewUpdate().Table("passwords")
	q.Set().Add("failed_attempts = 0")
	q.Set().Append("updated_at = now()")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("user_id", "password_hash", "failed_attempts", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&password.UserID, &password.PasswordHash, &password.FailedAttempts, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return password, nil
}
