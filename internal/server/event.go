package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/logger"
	"github.com/HardDie/event_tracker/internal/service"
	"github.com/HardDie/event_tracker/internal/utils"
)

type Event struct {
	service service.IEvent
}

func NewEvent(service service.IEvent) *Event {
	return &Event{
		service: service,
	}
}

func (s *Event) RegisterPublicRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	eventRouter := router.PathPrefix("").Subrouter()
	eventRouter.Use(middleware...)
}
func (s *Event) RegisterPrivateRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	eventRouter := router.PathPrefix("").Subrouter()
	eventRouter.HandleFunc("", s.CreateEvent).Methods(http.MethodPost)
	eventRouter.HandleFunc("/list", s.ListEvent).Methods(http.MethodPost)

	eventTypeRouter := eventRouter.PathPrefix("/types").Subrouter()
	eventTypeRouter.HandleFunc("", s.CreateEventType).Methods(http.MethodPost)
	eventTypeRouter.HandleFunc("", s.ListEventType).Methods(http.MethodGet)
	eventTypeRouter.HandleFunc("/{id:[0-9]+}", s.EditEventType).Methods(http.MethodPut)

	eventRouter.Use(middleware...)
}

/*
 * Private
 */

// swagger:parameters CreateEventTypeRequest
type CreateEventTypeRequest struct {
	// In: body
	Body struct {
		dto.CreateEventTypeDTO
	}
}

// swagger:response CreateEventTypeResponse
type CreateEventTypeResponse struct {
	// In: body
	Body struct {
		Data *entity.EventType `json:"data"`
	}
}

// swagger:route POST /api/v1/events/types EventType CreateEventTypeRequest
//
// # Create an event type
//
//	Responses:
//	  200: CreateEventTypeResponse
func (s *Event) CreateEventType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.CreateEventTypeDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventType, err := s.service.CreateType(ctx, userID, req)
	if err != nil {
		logger.Error.Println("Can't create event type:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = utils.Response(w, eventType)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

// swagger:parameters ListEventTypeRequest
type ListEventTypeRequest struct {
}

// swagger:response ListEventTypeResponse
type ListEventTypeResponse struct {
	// In: body
	Body struct {
		Data []*entity.EventType `json:"data"`
		Meta *utils.Meta         `json:"meta"`
	}
}

// swagger:route GET /api/v1/events/types EventType ListEventTypeRequest
//
// # Getting a list of all types of events
//
//	Responses:
//	  200: ListEventTypeResponse
func (s *Event) ListEventType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	eventTypes, total, err := s.service.ListType(ctx, userID)
	if err != nil {
		logger.Error.Println("Can't get event type list:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if eventTypes == nil {
		eventTypes = make([]*entity.EventType, 0)
	}

	meta := &utils.Meta{
		Total: total,
	}
	err = utils.ResponseWithMeta(w, eventTypes, meta)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

// swagger:parameters EditEventTypeRequest
type EditEventTypeRequest struct {
	// In: path
	ID int32 `json:"id"`
	// In: body
	Body struct {
		dto.EditEventTypeDTO
	}
}

// swagger:response EditEventTypeResponse
type EditEventTypeResponse struct {
	// In: body
	Body struct {
		Data *entity.EventType `json:"data"`
	}
}

// swagger:route PUT /api/v1/events/types/{id} EventType EditEventTypeRequest
//
// # Editing the event type
//
//	Responses:
//	  200: EditEventTypeResponse
func (s *Event) EditEventType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.EditEventTypeDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}

	req.ID, err = utils.GetInt32FromPath(r, "id")
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Bad id in path", http.StatusBadRequest)
		return
	}

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventType, err := s.service.EditType(ctx, userID, req)
	if err != nil {
		logger.Error.Println("Can't edit event type:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = utils.Response(w, eventType)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

// swagger:parameters CreateEventRequest
type CreateEventRequest struct {
	// In: body
	Body struct {
		dto.CreateEventDTO
	}
}

// swagger:response CreateEventResponse
type CreateEventResponse struct {
	// In: body
	Body struct {
		Data *entity.Event `json:"data"`
	}
}

// swagger:route POST /api/v1/events Event CreateEventRequest
//
// # Create an event
//
//	Responses:
//	  200: CreateEventResponse
func (s *Event) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.CreateEventDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := s.service.CreateEvent(ctx, userID, req)
	if err != nil {
		logger.Error.Println("Can't create event:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = utils.Response(w, event)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

// swagger:parameters ListEventRequest
type ListEventRequest struct {
	// In: body
	Body struct {
		dto.ListEventDTO
	}
}

// swagger:response ListEventResponse
type ListEventResponse struct {
	// In: body
	Body struct {
		Data []*entity.Event `json:"data"`
		Meta *utils.Meta     `json:"meta"`
	}
}

// swagger:route POST /api/v1/events/list Event ListEventRequest
//
// # Getting a list of events
//
//	Responses:
//	  200: ListEventResponse
func (s *Event) ListEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.ListEventDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, total, err := s.service.ListEvent(ctx, userID, req)
	if err != nil {
		logger.Error.Println("Can't get event list:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if events == nil {
		events = make([]*entity.Event, 0)
	}

	meta := &utils.Meta{
		Total: total,
	}
	err = utils.ResponseWithMeta(w, events, meta)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}
