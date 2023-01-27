package application

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/HardDie/event_tracker/internal/config"
	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/logger"
	"github.com/HardDie/event_tracker/internal/middleware"
	"github.com/HardDie/event_tracker/internal/migration"
	"github.com/HardDie/event_tracker/internal/repository"
	"github.com/HardDie/event_tracker/internal/server"
	"github.com/HardDie/event_tracker/internal/service"
)

type Application struct {
	Cfg    *config.Config
	DB     *db.DB
	Router *mux.Router
}

func Get() (*Application, error) {
	app := &Application{
		Cfg:    config.Get(),
		Router: mux.NewRouter(),
	}

	// Init DB
	newDB, err := db.Get(app.Cfg.DBPath)
	if err != nil {
		return nil, err
	}
	app.DB = newDB

	// Init migrations
	err = migration.NewMigrate(app.DB).Up()
	if err != nil {
		return nil, err
	}

	// Prepare router
	apiRouter := app.Router.PathPrefix("/api").Subrouter()
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()

	// Init repositories
	userRepository := repository.NewUser(app.DB)
	passwordRepository := repository.NewPassword(app.DB)
	sessionRepository := repository.NewSession(app.DB)
	eventRepository := repository.NewEvent(app.DB)
	friendRepository := repository.NewFriend(app.DB)

	// Init services
	systemService := service.NewSystem()
	authService := service.NewAuth(app.Cfg, userRepository, passwordRepository, sessionRepository)
	eventService := service.NewEvent(eventRepository)
	friendService := service.NewFriend(friendRepository, userRepository)

	// Init severs
	systemServer := server.NewSystem(systemService)
	authServer := server.NewAuth(app.Cfg, authService)
	userServer := server.NewUser(
		service.NewUser(userRepository, passwordRepository),
	)
	eventServer := server.NewEvent(eventService)
	friendServer := server.NewFriend(friendService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)
	timeoutMiddleware := middleware.NewTimeoutRequestMiddleware(time.Duration(app.Cfg.RequestTimeout) * time.Second)

	// Register servers
	systemRouter := v1Router.PathPrefix("/system").Subrouter()
	systemServer.RegisterPublicRouter(systemRouter, middleware.CorsMiddleware, timeoutMiddleware.RequestMiddleware)

	authRouter := v1Router.PathPrefix("/auth").Subrouter()
	authServer.RegisterPublicRouter(authRouter)
	authServer.RegisterPrivateRouter(authRouter, timeoutMiddleware.RequestMiddleware, authMiddleware.RequestMiddleware)

	userRouter := v1Router.PathPrefix("/user").Subrouter()
	userServer.RegisterPrivateRouter(userRouter, timeoutMiddleware.RequestMiddleware, authMiddleware.RequestMiddleware)

	eventRouter := v1Router.PathPrefix("/events").Subrouter()
	eventServer.RegisterPrivateRouter(eventRouter, timeoutMiddleware.RequestMiddleware, authMiddleware.RequestMiddleware)

	friendRouter := v1Router.PathPrefix("/friends").Subrouter()
	friendServer.RegisterPrivateRouter(friendRouter, timeoutMiddleware.RequestMiddleware, authMiddleware.RequestMiddleware)

	return app, nil
}

func (app *Application) Run() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		app.Stop()
		os.Exit(0)
	}()

	defer app.Stop()
	return http.ListenAndServe(app.Cfg.Port, app.Router)
}

func (app *Application) Stop() {
	err := app.DB.DB.Close()
	if err != nil {
		logger.Error.Println("error closing DB connection:", err.Error())
	}
	app.DB = nil
	log.Println("Done")
}
