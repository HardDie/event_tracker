package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/errs"
	"github.com/HardDie/event_tracker/internal/logger"
	"github.com/HardDie/event_tracker/internal/service"
	"github.com/HardDie/event_tracker/internal/utils"
)

type User struct {
	service service.IUser
}

func NewUser(service service.IUser) *User {
	return &User{
		service: service,
	}
}
func (s *User) RegisterPrivateRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	userRouter := router.PathPrefix("").Subrouter()
	userRouter.HandleFunc("/{id:[0-9]+}", s.Get).Methods(http.MethodGet)
	userRouter.HandleFunc("/password", s.Password).Methods(http.MethodPut)
	userRouter.HandleFunc("/profile", s.UpdateProfile).Methods(http.MethodPut)
	userRouter.HandleFunc("/image", s.UpdateImage).Methods(http.MethodPut)
	userRouter.Use(middleware...)
}

/*
 * Private
 */

// swagger:parameters UserGetRequest
type UserGetRequest struct {
	// In: path
	ID int32 `json:"id"`
}

// swagger:response UserGetResponse
type UserGetResponse struct {
	// In: body
	Body struct {
		Data *entity.User `json:"data"`
	}
}

// swagger:route GET /api/v1/user/{id} User UserGetRequest
//
// # Getting information about a user by ID
//
//	Responses:
//	  200: UserGetResponse
func (s *User) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	id, err := utils.GetInt32FromPath(r, "id")
	if err != nil {
		http.Error(w, "Bad id in path", http.StatusBadRequest)
		return
	}
	req := dto.GetUserDTO{
		ID: id,
	}

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.service.Get(ctx, id, userID)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	err = utils.Response(w, user)
	if err != nil {
		logger.Error.Println("error write to socket:", err.Error())
	}
}

// swagger:parameters UserPasswordRequest
type UserPasswordRequest struct {
	// In: body
	Body struct {
		dto.UpdatePasswordDTO
	}
}

// swagger:response UserPasswordResponse
type UserPasswordResponse struct {
}

// swagger:route PUT /api/v1/user/password User UserPasswordRequest
//
// # Updating the password for a user
//
//	Responses:
//	  200: UserPasswordResponse
func (s *User) Password(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.UpdatePasswordDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.service.Password(ctx, req, userID)
	if err != nil {
		errs.HttpError(w, err)
		return
	}
}

// swagger:parameters UserUpdateProfileRequest
type UserUpdateProfileRequest struct {
	// In: body
	Body struct {
		dto.UpdateProfileDTO
	}
}

// swagger:response UserUpdateProfileResponse
type UserUpdateProfileResponse struct {
	// In: body
	Body struct {
		Data *entity.User `json:"data"`
	}
}

// swagger:route PUT /api/v1/user/profile User UserUpdateProfileRequest
//
// # Updating user information
//
//	Responses:
//	  200: UserUpdateProfileResponse
func (s *User) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.UpdateProfileDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}
	req.ID = userID

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.service.UpdateProfile(ctx, req)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	err = utils.Response(w, user)
	if err != nil {
		logger.Error.Println("error write to socket:", err.Error())
	}
}

// swagger:parameters UserUpdateImageRequest
type UserUpdateImageRequest struct {
	// In: body
	Body struct {
		dto.UpdateProfileImageDTO
	}
}

// swagger:response UserUpdateImageResponse
type UserUpdateImageResponse struct {
	// In: body
	Body struct {
		Data *entity.User `json:"data"`
	}
}

// swagger:route PUT /api/v1/user/image User UserUpdateImageRequest
//
// # Updating user profile image
//
//	Responses:
//	  200: UserUpdateImageResponse
func (s *User) UpdateImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.UpdateProfileImageDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}
	req.ID = userID

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.service.UpdateImage(ctx, req)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	err = utils.Response(w, user)
	if err != nil {
		logger.Error.Println("error write to socket:", err.Error())
	}
}
