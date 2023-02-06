package service

import (
	"context"

	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/errs"
	"github.com/HardDie/event_tracker/internal/logger"
	"github.com/HardDie/event_tracker/internal/repository"
)

type IEvent interface {
	CreateType(ctx context.Context, userID int32, req *dto.CreateEventTypeDTO) (*entity.EventType, error)
	DeleteType(ctx context.Context, userID int32, id int32) error
	ListType(ctx context.Context, userID int32) ([]*entity.EventType, int32, error)
	EditType(ctx context.Context, userID int32, req *dto.EditEventTypeDTO) (*entity.EventType, error)

	CreateEvent(ctx context.Context, userID int32, req *dto.CreateEventDTO) (*entity.Event, error)
	DeleteEvent(ctx context.Context, userID int32, id int32) error
	ListEvent(ctx context.Context, userId int32, req *dto.ListEventDTO) ([]*entity.Event, int32, error)

	FriendsFeed(ctx context.Context, userID int32) ([]*dto.FeedResponseDTO, int32, error)
}

type Event struct {
	repository repository.IEvent

	db *db.DB
}

func NewEvent(db *db.DB, repository repository.IEvent) *Event {
	return &Event{
		db:         db,
		repository: repository,
	}
}

func (s *Event) CreateType(ctx context.Context, userID int32, req *dto.CreateEventTypeDTO) (*entity.EventType, error) {
	res, err := s.repository.CreateType(s.db.DB, ctx, userID, req.Name, req.IsVisible)
	if err != nil {
		logger.Error.Printf("error create type: %v", err.Error())
		return nil, errs.InternalError
	}
	return res, nil
}
func (s *Event) DeleteType(ctx context.Context, userID int32, id int32) error {
	err := s.repository.DeleteType(s.db.DB, ctx, userID, id)
	if err != nil {
		logger.Error.Printf("error delete type: %v", err.Error())
		return errs.InternalError
	}
	return nil
}
func (s *Event) ListType(ctx context.Context, userID int32) ([]*entity.EventType, int32, error) {
	res, cnt, err := s.repository.ListType(s.db.DB, ctx, userID, false)
	if err != nil {
		logger.Error.Printf("error list type: %v", err.Error())
		return nil, 0, errs.InternalError
	}
	return res, cnt, nil
}
func (s *Event) EditType(ctx context.Context, userID int32, req *dto.EditEventTypeDTO) (*entity.EventType, error) {
	res, err := s.repository.EditType(s.db.DB, ctx, userID, req.ID, req.Name, req.IsVisible)
	if err != nil {
		logger.Error.Printf("error edit type: %v", err.Error())
		return nil, errs.InternalError
	}
	return res, nil
}

func (s *Event) CreateEvent(ctx context.Context, userID int32, req *dto.CreateEventDTO) (*entity.Event, error) {
	res, err := s.repository.CreateEvent(s.db.DB, ctx, userID, req.EventTypeID, req.Date)
	if err != nil {
		logger.Error.Printf("error create event: %v", err.Error())
		return nil, errs.InternalError
	}
	return res, nil
}
func (s *Event) DeleteEvent(ctx context.Context, userID int32, id int32) error {
	err := s.repository.DeleteEvent(s.db.DB, ctx, userID, id)
	if err != nil {
		logger.Error.Printf("error delete event: %v", err.Error())
		return errs.InternalError
	}
	return nil
}
func (s *Event) ListEvent(ctx context.Context, userID int32, req *dto.ListEventDTO) ([]*entity.Event, int32, error) {
	reqUserID := userID
	if req.UserID != nil {
		reqUserID = *req.UserID
	}
	res, cnt, err := s.repository.ListEvent(s.db.DB, ctx, &dto.ListEventFilter{
		UserID:      reqUserID,
		TypeID:      req.TypeID,
		OnlyVisible: reqUserID != userID,
		PeriodType:  req.PeriodType,
		Date:        req.Date,
	})
	if err != nil {
		logger.Error.Printf("error list event: %v", err.Error())
		return nil, 0, errs.InternalError
	}
	return res, cnt, nil
}

func (s *Event) FriendsFeed(ctx context.Context, userID int32) ([]*dto.FeedResponseDTO, int32, error) {
	res, cnt, err := s.repository.FriendsFeed(s.db.DB, ctx, userID)
	if err != nil {
		logger.Error.Printf("error list event: %v", err.Error())
		return nil, 0, errs.InternalError
	}
	return res, cnt, nil
}
