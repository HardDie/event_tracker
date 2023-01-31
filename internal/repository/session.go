package repository

import (
	"context"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/event_tracker/internal/entity"
)

type ISession interface {
	CreateOrUpdate(tx IQuery, ctx context.Context, userID int32, sessionHash string) (*entity.Session, error)
	GetByUserID(tx IQuery, ctx context.Context, sessionHash string) (*entity.Session, error)
	DeleteByID(tx IQuery, ctx context.Context, id int32) error
}

type Session struct {
}

func NewSession() *Session {
	return &Session{}
}

func (r *Session) CreateOrUpdate(tx IQuery, ctx context.Context, userID int32, sessionHash string) (*entity.Session, error) {
	session := &entity.Session{
		UserID:      userID,
		SessionHash: sessionHash,
	}

	q := gosql.NewInsert().Into("sessions")
	q.Columns().Add("user_id", "session_hash")
	q.Columns().Arg(userID, sessionHash)
	q.Conflict().Object("user_id").Action("UPDATE").Set().
		Add("session_hash = EXCLUDED.session_hash", "updated_at = datetime('now')", "deleted_at = NULL")
	q.Returning().Add("id", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}
func (r *Session) GetByUserID(tx IQuery, ctx context.Context, sessionHash string) (*entity.Session, error) {
	session := &entity.Session{
		SessionHash: sessionHash,
	}

	q := gosql.NewSelect().From("sessions")
	q.Columns().Add("id", "user_id", "created_at", "updated_at")
	q.Where().AddExpression("session_hash = ?", sessionHash)
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}
func (r *Session) DeleteByID(tx IQuery, ctx context.Context, id int32) error {
	q := gosql.NewUpdate().Table("sessions")
	q.Set().Add("deleted_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("id")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
