package main

import (
	"github.com/saalikmubeen/goravel"
	"github.com/saalikmubeen/goravel-demo-app/handlers"
	"github.com/saalikmubeen/goravel-demo-app/middleware"
	"github.com/saalikmubeen/goravel-demo-app/models"
)

type application struct {
	App        *goravel.Goravel
	Handlers   *handlers.Handlers
	Models     *models.Models
	Middleware *middleware.Middleware
}

func main() {

	app := initGoravel()
	app.App.ListenAndServe()

}
