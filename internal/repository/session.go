package repository

import (
	"context"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/entity"
)

type ISession interface {
	CreateOrUpdate(ctx context.Context, userID int32, sessionHash string) (*entity.Session, error)
	GetByUserID(ctx context.Context, sessionHash string) (*entity.Session, error)
	DeleteByID(ctx context.Context, id int32) error
}

type Session struct {
	db *db.DB
}

func NewSession(db *db.DB) *Session {
	return &Session{
		db: db,
	}
}

func (r *Session) CreateOrUpdate(ctx context.Context, userID int32, sessionHash string) (*entity.Session, error) {
	session := &entity.Session{
		UserID:      userID,
		SessionHash: sessionHash,
	}

	row := r.db.DB.QueryRowContext(ctx, `
INSERT INTO sessions (user_id, session_hash)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
SET session_hash = $2, updated_at = datetime('now'), deleted_at = NULL
RETURNING id, created_at, updated_at`, userID, sessionHash)

	err := row.Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}
func (r *Session) GetByUserID(ctx context.Context, sessionHash string) (*entity.Session, error) {
	session := &entity.Session{
		SessionHash: sessionHash,
	}

	q := gosql.NewSelect().From("sessions")
	q.Columns().Add("id", "user_id", "created_at", "updated_at")
	q.Where().AddExpression("session_hash = ?", sessionHash)
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}
func (r *Session) DeleteByID(ctx context.Context, id int32) error {
	q := gosql.NewUpdate().Table("sessions")
	q.Set().Add("deleted_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("id")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
