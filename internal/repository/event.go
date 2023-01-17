package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/logger"
)

type IEvent interface {
	CreateType(ctx context.Context, userID int32, name string, isVisible bool) (*entity.EventType, error)
	DeleteType(ctx context.Context, userID, id int32) error
	ListType(ctx context.Context, userID int32, onlyVisible bool) ([]*entity.EventType, int32, error)
	EditType(ctx context.Context, userID, id int32, name string, isVisible bool) (*entity.EventType, error)

	CreateEvent(ctx context.Context, userID, eventTypeID int32, date time.Time) (*entity.Event, error)
	DeleteEvent(ctx context.Context, userID, id int32) error
	ListEvent(ctx context.Context, filter *dto.ListEventFilter) ([]*entity.Event, int32, error)
}
type Event struct {
	db *db.DB
}

func NewEvent(db *db.DB) *Event {
	return &Event{
		db: db,
	}
}

func (r *Event) CreateType(ctx context.Context, userID int32, name string, isVisible bool) (*entity.EventType, error) {
	eventType := &entity.EventType{
		EventType: name,
		IsVisible: isVisible,
	}

	q := gosql.NewInsert().Into("event_types")
	q.Columns().Add("event_type", "is_visible", "user_id")
	q.Columns().Arg(name, isVisible, userID)
	q.Returning().Add("id", "created_at", "updated_at")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&eventType.ID, &eventType.CreatedAt, &eventType.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorEventTypeNotExist
		}
		logger.Error.Println(err.Error())
		return nil, ErrorInternal
	}
	return eventType, nil
}
func (r *Event) DeleteType(ctx context.Context, userID, id int32) error {
	q := gosql.NewUpdate().Table("event_types")
	q.Set().Add("deleted_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("id")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorEventTypeNotExist
		}
		logger.Error.Println(err.Error())
		return ErrorInternal
	}
	return nil
}
func (r *Event) ListType(ctx context.Context, userID int32, onlyVisible bool) ([]*entity.EventType, int32, error) {
	var res []*entity.EventType

	q := gosql.NewSelect().From("event_types")
	q.Columns().Add("id", "user_id", "event_type", "is_visible", "created_at", "updated_at")
	q.Where().AddExpression("deleted_at IS NULL")
	q.Where().AddExpression("user_id = ?", userID)
	if onlyVisible {
		q.Where().AddExpression("is_visible")
	}
	q.AddOrder("event_type")
	rows, err := r.db.DB.QueryContext(ctx, q.String(), q.GetArguments()...)
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, 0, ErrorInternal
	}
	defer rows.Close()

	for rows.Next() {
		eventType := &entity.EventType{}
		err = rows.Scan(&eventType.ID, &eventType.UserID, &eventType.EventType, &eventType.IsVisible, &eventType.CreatedAt, &eventType.UpdatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, 0, ErrorEventTypeNotExist
			}
			logger.Error.Println(err.Error())
			return nil, 0, ErrorInternal
		}
		res = append(res, eventType)
	}

	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, ErrorEventTypeNotExist
		}
		logger.Error.Println(err.Error())
		return nil, 0, ErrorInternal
	}

	return res, int32(len(res)), nil
}
func (r *Event) EditType(ctx context.Context, userID, id int32, name string, isVisible bool) (*entity.EventType, error) {
	eventType := &entity.EventType{
		ID:        id,
		UserID:    userID,
		EventType: name,
		IsVisible: isVisible,
	}

	q := gosql.NewUpdate().Table("event_types")
	q.Set().Append("event_type = ?", name)
	q.Set().Append("is_visible = ?", isVisible)
	q.Set().Append("updated_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("created_at", "updated_at")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&eventType.CreatedAt, &eventType.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorEventTypeNotExist
		}
		logger.Error.Println(err.Error())
		return nil, ErrorInternal
	}
	return eventType, nil
}

func (r *Event) CreateEvent(ctx context.Context, userID, typeID int32, date time.Time) (*entity.Event, error) {
	date = timeToYMD(date)
	event := &entity.Event{
		UserID: userID,
		TypeID: typeID,
		Date:   date,
	}

	q := gosql.NewInsert().Into("events")
	q.Columns().Add("user_id", "type_id", "date")
	q.Columns().Arg(userID, typeID, date)
	q.Returning().Add("id", "created_at", "updated_at")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorEventNotExist
		}
		logger.Error.Println(err.Error())
		return nil, ErrorInternal
	}
	return event, nil
}
func (r *Event) DeleteEvent(ctx context.Context, userID, id int32) error {
	q := gosql.NewUpdate().Table("events")
	q.Set().Add("deleted_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("id")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorEventNotExist
		}
		logger.Error.Println(err.Error())
		return ErrorInternal
	}
	return nil
}
func (r *Event) ListEvent(ctx context.Context, filter *dto.ListEventFilter) ([]*entity.Event, int32, error) {
	var res []*entity.Event

	q := gosql.NewSelect().From("events ev")
	q.Relate("JOIN event_types et ON ev.type_id = et.id")
	q.Columns().Add("ev.id", "ev.user_id", "ev.type_id", "ev.date", "ev.created_at", "ev.updated_at")
	q.Where().AddExpression("ev.deleted_at IS NULL")
	q.Where().AddExpression("et.deleted_at IS NULL")
	q.Where().AddExpression("ev.user_id = ?", filter.UserID)
	if filter.TypeID != nil {
		q.Where().AddExpression("ev.type_id = ?", filter.TypeID)
	}
	if filter.OnlyVisible {
		q.Where().AddExpression("et.is_visible")
	}
	switch filter.PeriodType {
	case dto.PeriodDay:
		q.Where().AddExpression("ev.date = ?", timeToYMD(filter.Date))
	case dto.PeriodMonth:
		first, last := timeToYM(filter.Date)
		q.Where().AddExpression("ev.date >= ?", first)
		q.Where().AddExpression("ev.date <= ?", last)
	case dto.PeriodYear:
		first, last := timeToY(filter.Date)
		q.Where().AddExpression("ev.date >= ?", first)
		q.Where().AddExpression("ev.date <= ?", last)
	}
	q.AddOrder("ev.created_at")

	rows, err := r.db.DB.QueryContext(ctx, q.String(), q.GetArguments()...)
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, 0, ErrorInternal
	}
	defer rows.Close()

	for rows.Next() {
		event := &entity.Event{}
		err = rows.Scan(&event.ID, &event.UserID, &event.TypeID, &event.Date, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, 0, ErrorEventNotExist
			}
			logger.Error.Println(err.Error())
			return nil, 0, ErrorInternal
		}
		res = append(res, event)
	}

	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, ErrorEventNotExist
		}
		logger.Error.Println(err.Error())
		return nil, 0, ErrorInternal
	}

	return res, int32(len(res)), nil
}

func timeToYMD(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}
func timeToYM(date time.Time) (time.Time, time.Time) {
	firstDayOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
	return firstDayOfMonth, lastDayOfMonth
}
func timeToY(date time.Time) (time.Time, time.Time) {
	return time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(date.Year(), 12, 31, 0, 0, 0, 0, time.UTC)
}