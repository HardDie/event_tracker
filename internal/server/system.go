package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/event_tracker/internal/logger"
	"github.com/HardDie/event_tracker/internal/service"
)

type System struct {
	service service.ISystem
}

func NewSystem(service service.ISystem) *System {
	return &System{
		service: service,
	}
}

func (s *System) RegisterPublicRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	systemRouter := router.PathPrefix("").Subrouter()
	systemRouter.HandleFunc("/swagger", s.Swagger).Methods(http.MethodGet)
	systemRouter.Use(middleware...)
}

/*
 * Public
 */

// swagger:parameters SwaggerRequest
type SwaggerRequest struct {
}

// swagger:response SwaggerResponse
type SwaggerResponse struct {
	// In: body
	Body []byte
}

// swagger:route GET /api/v1/system/swagger System SwaggerRequest
//
// # Get the yaml-file of the swagger description
//
//	Responses:
//	  200: SwaggerResponse
func (s *System) Swagger(w http.ResponseWriter, r *http.Request) {
	data, err := s.service.GetSwagger()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(data)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}
