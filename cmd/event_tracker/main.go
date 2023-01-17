// Package main Event Tracker.
//
// Entry point for the application.
//
// Terms Of Service:
//
//	Schemes: http
//	Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- token:
//
//	SecurityDefinitions:
//	token:
//	     type: apiKey
//	     name: Cookie
//	     in: header
//
// swagger:meta
package main

import (
	"github.com/HardDie/event_tracker/internal/application"
	"github.com/HardDie/event_tracker/internal/logger"
)

func main() {
	app, err := application.Get()
	if err != nil {
		logger.Error.Fatal(err)
	}
	logger.Info.Println("Server listen on", app.Cfg.Port)
	err = app.Run()
	if err != nil {
		logger.Error.Fatal(err)
	}
}
