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

type Friend struct {
	service service.IFriend
}

func NewFriend(service service.IFriend) *Friend {
	return &Friend{
		service: service,
	}
}

func (s *Friend) RegisterPrivateRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	friendRouter := router.PathPrefix("").Subrouter()
	friendRouter.HandleFunc("", s.FriendList).Methods(http.MethodGet)
	friendRouter.HandleFunc("/invites", s.InviteFriend).Methods(http.MethodPost)
	friendRouter.HandleFunc("/invites", s.InviteList).Methods(http.MethodGet)
	friendRouter.HandleFunc("/invites/{id:[0-9]+}", s.InviteAccept).Methods(http.MethodPost)
	friendRouter.HandleFunc("/invites/{id:[0-9]+}", s.InviteReject).Methods(http.MethodDelete)
	friendRouter.Use(middleware...)
}

/*
 * Private
 */

// swagger:parameters InviteFriendRequest
type InviteFriendRequest struct {
	// In: body
	Body struct {
		dto.InviteFriendDTO
	}
}

// swagger:response InviteFriendResponse
type InviteFriendResponse struct {
}

// swagger:route POST /api/v1/friends/invites Friend InviteFriendRequest
//
// # Invite user to friends
//
//	Responses:
//	  200: InviteFriendResponse
func (s *Friend) InviteFriend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.InviteFriendDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}
	req.ID = userID

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.service.InviteFriend(ctx, req)
	if err != nil {
		logger.Error.Println("Can't creation invitation:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

// swagger:parameters InviteListRequest
type InviteListRequest struct {
}

// swagger:response InviteListResponse
type InviteListResponse struct {
	// In: body
	Body struct {
		Data []*dto.InviteListResponseDTO `json:"data"`
		Meta *utils.Meta                  `json:"meta"`
	}
}

// swagger:route GET /api/v1/friends/invites Friend InviteListRequest
//
// # Get a list of pending invitations
//
//	Responses:
//	  200: InviteListResponse
func (s *Friend) InviteList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	invites, total, err := s.service.InviteListPending(ctx, userID)
	if err != nil {
		logger.Error.Println("Can't get list of invitations:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if invites == nil {
		invites = make([]*dto.InviteListResponseDTO, 0)
	}

	err = utils.ResponseWithMeta(w, invites, &utils.Meta{
		Total: total,
	})
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

// swagger:parameters InviteAcceptRequest
type InviteAcceptRequest struct {
	// In:path
	ID int32 `json:"id"`
}

// swagger:response InviteAcceptResponse
type InviteAcceptResponse struct {
}

// swagger:route POST /api/v1/friends/invites/{id} Friend InviteAcceptRequest
//
// # Accept a friend request
//
//	Responses:
//	  200: InviteAcceptResponse
func (s *Friend) InviteAccept(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	id, err := utils.GetInt32FromPath(r, "id")
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Bad id in path", http.StatusBadRequest)
		return
	}

	err = s.service.AcceptFriendship(ctx, userID, id)
	if err != nil {
		logger.Error.Println("Can't accept friendship:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

// swagger:parameters InviteRejectRequest
type InviteRejectRequest struct {
	// In:path
	ID int32 `json:"id"`
}

// swagger:response InviteRejectResponse
type InviteRejectResponse struct {
}

// swagger:route DELETE /api/v1/friends/invites/{id} Friend InviteRejectRequest
//
// # Reject a friend request
//
//	Responses:
//	  200: InviteRejectResponse
func (s *Friend) InviteReject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	id, err := utils.GetInt32FromPath(r, "id")
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Bad id in path", http.StatusBadRequest)
		return
	}

	err = s.service.RejectFriendship(ctx, userID, id)
	if err != nil {
		logger.Error.Println("Can't reject friendship:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

// swagger:parameters FriendListRequest
type FriendListRequest struct {
}

// swagger:response FriendListResponse
type FriendListResponse struct {
	// In: body
	Body struct {
		Data []*entity.User `json:"data"`
		Meta *utils.Meta    `json:"meta"`
	}
}

// swagger:route GET /api/v1/friends Friend FriendListRequest
//
// # Get a list of friends
//
//	Responses:
//	  200: FriendListResponse
func (s *Friend) FriendList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	friends, total, err := s.service.ListOfFriends(ctx, userID)
	if err != nil {
		logger.Error.Println("Can't get list of friends:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if friends == nil {
		friends = make([]*entity.User, 0)
	}

	err = utils.ResponseWithMeta(w, friends, &utils.Meta{
		Total: total,
	})
	if err != nil {
		logger.Error.Println(err.Error())
	}
}
