package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/event_tracker/internal/config"
	"github.com/HardDie/event_tracker/internal/dto"
	"github.com/HardDie/event_tracker/internal/entity"
	"github.com/HardDie/event_tracker/internal/logger"
	"github.com/HardDie/event_tracker/internal/service"
	"github.com/HardDie/event_tracker/internal/utils"
)

type Auth struct {
	service service.IAuth
	cfg     *config.Config
}

func NewAuth(cfg *config.Config, service service.IAuth) *Auth {
	return &Auth{
		cfg:     cfg,
		service: service,
	}
}
func (s *Auth) RegisterPublicRouter(router *mux.Router) {
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.HandleFunc("/register", s.Register).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", s.Login).Methods(http.MethodPost)
}
func (s *Auth) RegisterPrivateRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.HandleFunc("/user", s.User).Methods(http.MethodGet)
	authRouter.HandleFunc("/logout", s.Logout).Methods(http.MethodPost)
	authRouter.Use(middleware...)
}

/*
 * Public
 */

// swagger:parameters AuthRegisterRequest
type AuthRegisterRequest struct {
	// In: body
	Body struct {
		dto.RegisterDTO
	}
}

// swagger:response AuthRegisterResponse
type AuthRegisterResponse struct {
	// In: body
	Body struct {
		*entity.Session
	}
}

// swagger:route POST /api/v1/auth/register Auth AuthRegisterRequest
//
// # Registration form
//
//	Responses:
//	  200: AuthRegisterResponse
func (s *Auth) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &dto.RegisterDTO{}
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

	user, err := s.service.Register(ctx, req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	session, err := s.service.GenerateCookie(ctx, user.ID)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	utils.SetSessionCookie(session.Session, w)
	err = utils.Response(w, session)
	if err != nil {
		logger.Error.Printf(err.Error())
	}
}

// swagger:parameters AuthLoginRequest
type AuthLoginRequest struct {
	// In: body
	Body struct {
		dto.LoginDTO
	}
}

// swagger:response AuthLoginResponse
type AuthLoginResponse struct {
	// In: body
	Body struct {
		*entity.Session
	}
}

// swagger:route POST /api/v1/auth/login Auth AuthLoginRequest
//
// # Login form
//
//	Responses:
//	  200: AuthLoginResponse
func (s *Auth) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &dto.LoginDTO{}
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

	user, err := s.service.Login(ctx, req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	session, err := s.service.GenerateCookie(ctx, user.ID)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	utils.SetSessionCookie(session.Session, w)
	err = utils.Response(w, session)
	if err != nil {
		logger.Error.Printf(err.Error())
	}
}

/*
 * Private
 */

// swagger:parameters AuthUserRequest
type AuthUserRequest struct {
}

// swagger:response AuthUserResponse
type AuthUserResponse struct {
	// In: body
	Body struct {
		Data *entity.User `json:"data"`
	}
}

// swagger:route GET /api/v1/auth/user Auth AuthUserRequest
//
// # Getting information about the current user
//
//	Responses:
//	  200: AuthUserResponse
func (s *Auth) User(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	user, err := s.service.GetUserInfo(ctx, userID)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	err = utils.Response(w, user)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

// swagger:parameters AuthLogoutRequest
type AuthLogoutRequest struct {
}

// swagger:response AuthLogoutResponse
type AuthLogoutResponse struct {
}

// swagger:route POST /api/v1/auth/logout Auth AuthLogoutRequest
//
// # Close the current session
//
//	Responses:
//	  200: AuthLogoutResponse
func (s *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := utils.GetSessionFromContext(ctx)

	err := s.service.Logout(ctx, session.ID)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	utils.DeleteSessionCookie(w)
}
