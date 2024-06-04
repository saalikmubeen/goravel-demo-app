package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) ApiRoutes() http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(mux chi.Router) {
		// add any api routes here

		// User routes
		r.Get("/users", app.Handlers.GetAllUsers)
		r.Post("/users", app.Handlers.CreateUser)
		r.Get("/users/{id}", app.Handlers.GetUser)
		r.Put("/users/{id}", app.Handlers.UpdateUser)
		r.Delete("/users/{id}", app.Handlers.DeleteUser)

		// Cache routes
		r.Post("/save-in-cache", app.Handlers.SaveInCache)
		r.Post("/get-from-cache", app.Handlers.GetFromCache)
		r.Delete("/delete-from-cache", app.Handlers.DeleteFromCache)
		r.Delete("/empty-cache", app.Handlers.EmptyCache)
	})

	return r
}
