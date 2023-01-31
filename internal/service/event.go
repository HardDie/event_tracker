package service

import (
	"context"

	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
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
	return s.repository.CreateType(s.db.DB, ctx, userID, req.Name, req.IsVisible)
}
func (s *Event) DeleteType(ctx context.Context, userID int32, id int32) error {
	return s.repository.DeleteType(s.db.DB, ctx, userID, id)
}
func (s *Event) ListType(ctx context.Context, userID int32) ([]*entity.EventType, int32, error) {
	return s.repository.ListType(s.db.DB, ctx, userID, false)
}
func (s *Event) EditType(ctx context.Context, userID int32, req *dto.EditEventTypeDTO) (*entity.EventType, error) {
	return s.repository.EditType(s.db.DB, ctx, userID, req.ID, req.Name, req.IsVisible)
}

func (s *Event) CreateEvent(ctx context.Context, userID int32, req *dto.CreateEventDTO) (*entity.Event, error) {
	return s.repository.CreateEvent(s.db.DB, ctx, userID, req.EventTypeID, req.Date)
}
func (s *Event) DeleteEvent(ctx context.Context, userID int32, id int32) error {
	return s.repository.DeleteEvent(s.db.DB, ctx, userID, id)
}
func (s *Event) ListEvent(ctx context.Context, userID int32, req *dto.ListEventDTO) ([]*entity.Event, int32, error) {
	reqUserID := userID
	if req.UserID != nil {
		reqUserID = *req.UserID
	}
	return s.repository.ListEvent(s.db.DB, ctx, &dto.ListEventFilter{
		UserID:      reqUserID,
		TypeID:      req.TypeID,
		OnlyVisible: reqUserID != userID,
		PeriodType:  req.PeriodType,
		Date:        req.Date,
	})
}

func (s *Event) FriendsFeed(ctx context.Context, userID int32) ([]*dto.FeedResponseDTO, int32, error) {
	return s.repository.FriendsFeed(s.db.DB, ctx, userID)
}
