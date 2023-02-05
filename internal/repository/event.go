package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/HardDie/godb/v2"
	"github.com/dimonrus/gosql"

	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/logger"
)

type IEvent interface {
	CreateType(tx godb.Queryer, ctx context.Context, userID int32, name string, isVisible bool) (*entity.EventType, error)
	DeleteType(tx godb.Queryer, ctx context.Context, userID, id int32) error
	ListType(tx godb.Queryer, ctx context.Context, userID int32, onlyVisible bool) ([]*entity.EventType, int32, error)
	EditType(tx godb.Queryer, ctx context.Context, userID, id int32, name string, isVisible bool) (*entity.EventType, error)

	CreateEvent(tx godb.Queryer, ctx context.Context, userID, eventTypeID int32, date time.Time) (*entity.Event, error)
	DeleteEvent(tx godb.Queryer, ctx context.Context, userID, id int32) error
	ListEvent(tx godb.Queryer, ctx context.Context, filter *dto.ListEventFilter) ([]*entity.Event, int32, error)

	FriendsFeed(tx godb.Queryer, ctx context.Context, userID int32) ([]*dto.FeedResponseDTO, int32, error)
}
type Event struct {
}

func NewEvent() *Event {
	return &Event{}
}

func (r *Event) CreateType(tx godb.Queryer, ctx context.Context, userID int32, name string, isVisible bool) (*entity.EventType, error) {
	eventType := &entity.EventType{
		EventType: name,
		IsVisible: isVisible,
	}

	q := gosql.NewInsert().Into("event_types")
	q.Columns().Add("event_type", "is_visible", "user_id")
	q.Columns().Arg(name, isVisible, userID)
	q.Returning().Add("id", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

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
func (r *Event) DeleteType(tx godb.Queryer, ctx context.Context, userID, id int32) error {
	q := gosql.NewUpdate().Table("event_types")
	q.Set().Add("deleted_at = now()")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("id")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

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
func (r *Event) ListType(tx godb.Queryer, ctx context.Context, userID int32, onlyVisible bool) ([]*entity.EventType, int32, error) {
	var res []*entity.EventType

	q := gosql.NewSelect().From("event_types")
	q.Columns().Add("id", "user_id", "event_type", "is_visible", "created_at", "updated_at")
	q.Where().AddExpression("deleted_at IS NULL")
	q.Where().AddExpression("user_id = ?", userID)
	if onlyVisible {
		q.Where().AddExpression("is_visible")
	}
	q.AddOrder("event_type")
	rows, err := tx.QueryContext(ctx, q.String(), q.GetArguments()...)
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
func (r *Event) EditType(tx godb.Queryer, ctx context.Context, userID, id int32, name string, isVisible bool) (*entity.EventType, error) {
	eventType := &entity.EventType{
		ID:        id,
		UserID:    userID,
		EventType: name,
		IsVisible: isVisible,
	}

	q := gosql.NewUpdate().Table("event_types")
	q.Set().Append("event_type = ?", name)
	q.Set().Append("is_visible = ?", isVisible)
	q.Set().Append("updated_at = now()")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

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

func (r *Event) CreateEvent(tx godb.Queryer, ctx context.Context, userID, typeID int32, date time.Time) (*entity.Event, error) {
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
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

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
func (r *Event) DeleteEvent(tx godb.Queryer, ctx context.Context, userID, id int32) error {
	q := gosql.NewUpdate().Table("events")
	q.Set().Add("deleted_at = now()")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("id")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

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
func (r *Event) ListEvent(tx godb.Queryer, ctx context.Context, filter *dto.ListEventFilter) ([]*entity.Event, int32, error) {
	var res []*entity.Event

	q := gosql.NewSelect().From("events e")
	q.Relate("JOIN event_types et ON e.type_id = et.id")
	q.Columns().Add("e.id", "e.user_id", "e.type_id", "e.date", "e.created_at", "e.updated_at")
	q.Where().AddExpression("e.deleted_at IS NULL")
	q.Where().AddExpression("et.deleted_at IS NULL")
	q.Where().AddExpression("e.user_id = ?", filter.UserID)
	if filter.TypeID != nil {
		q.Where().AddExpression("e.type_id = ?", filter.TypeID)
	}
	if filter.OnlyVisible {
		q.Where().AddExpression("et.is_visible")
	}
	switch filter.PeriodType {
	case dto.PeriodDay:
		q.Where().AddExpression("e.date = ?", timeToYMD(filter.Date))
	case dto.PeriodMonth:
		first, last := timeToYM(filter.Date)
		q.Where().AddExpression("e.date >= ?", first)
		q.Where().AddExpression("e.date <= ?", last)
	case dto.PeriodYear:
		first, last := timeToY(filter.Date)
		q.Where().AddExpression("e.date >= ?", first)
		q.Where().AddExpression("e.date <= ?", last)
	}
	q.AddOrder("e.created_at")

	rows, err := tx.QueryContext(ctx, q.String(), q.GetArguments()...)
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

func (r *Event) FriendsFeed(tx godb.Queryer, ctx context.Context, userID int32) ([]*dto.FeedResponseDTO, int32, error) {
	var res []*dto.FeedResponseDTO

	q := gosql.NewSelect().From("friends f")
	q.Columns().Add("e.id", "u.id", "et.id", "et.event_type", "e.date", "e.created_at")
	q.Relate("JOIN users u ON f.with_user_id = u.id")
	q.Relate("JOIN events e ON f.with_user_id = e.user_id")
	q.Relate("JOIN event_types et ON e.type_id = et.id")
	q.Where().AddExpression("f.user_id = ?", userID)
	q.Where().AddExpression("f.deleted_at IS NULL")
	q.Where().AddExpression("u.deleted_at IS NULL")
	q.Where().AddExpression("et.deleted_at IS NULL")
	q.Where().AddExpression("et.is_visible")
	q.AddOrder("e.id DESC")
	q.SetPagination(100, 0)

	rows, err := tx.QueryContext(ctx, q.String(), q.GetArguments()...)
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, 0, ErrorInternal
	}
	defer rows.Close()

	for rows.Next() {
		event := &dto.FeedResponseDTO{}
		err = rows.Scan(&event.EventID, &event.UserID, &event.EventTypeID, &event.EventType, &event.Date, &event.CreatedAt)
		if err != nil {
			logger.Error.Println(err.Error())
			return nil, 0, ErrorInternal
		}
		res = append(res, event)
	}

	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, nil
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
