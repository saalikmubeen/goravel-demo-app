package main

import (
	"log"
	"os"

	"github.com/saalikmubeen/goravel"
	"github.com/saalikmubeen/goravel-demo-app/handlers"
	"github.com/saalikmubeen/goravel-demo-app/middleware"
	"github.com/saalikmubeen/goravel-demo-app/models"
)

func initGoravel() *application {

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init goravel struct
	gor := &goravel.Goravel{}
	err = gor.New(path)
	if err != nil {
		log.Fatal(err)
	}

	gor.AppName = "goravel-demo-app"

	gor.InfoLog.Println("Debug is set to", gor.Debug)

	// Initialize models
	models := models.New(gor.DB)

	// Initialize handlers
	handlers := &handlers.Handlers{
		App:    gor,
		Models: models,
	}

	// initialize middleware
	midleware := &middleware.Middleware{
		App:    gor,
		Models: models,
	}

	app := &application{
		App:        gor,
		Handlers:   handlers,
		Models:     models,
		Middleware: midleware,
	}

	// initialize routes
	// app.App.Routes = app.routes()
	app.routes()

	return app
}
