package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {

	// Add your middleware here
	app.App.Routes.Use(app.Middleware.CheckRememberMe)

	// Add your routes here
	app.App.Routes.Get("/", app.Handlers.Home)
	app.App.Routes.Get("/go-page", app.Handlers.GoPage)
	app.App.Routes.Get("/jet-page", app.Handlers.JetPage)
	app.App.Routes.Get("/sessions", app.Handlers.SessionTest)
	app.App.Routes.Get("/json", app.Handlers.JSON)
	app.App.Routes.Get("/xml", app.Handlers.XML)
	app.App.Routes.Get("/download-file", app.Handlers.DownloadFile)
	app.App.Routes.Get("/crypto", app.Handlers.TestCrypto)

	app.App.Routes.Get("/cache-test", app.Handlers.ShowCachePage)

	// Auth routes
	app.App.Routes.Get("/users/login", app.Handlers.UserLogin)
	app.App.Routes.Post("/users/login", app.Handlers.PostUserLogin)
	app.App.Routes.Get("/users/signup", app.Handlers.PostUserSignup)
	app.App.Routes.Post("/users/signup", app.Handlers.PostUserSignup)
	app.App.Routes.Get("/users/logout", app.Handlers.Logout)
	app.App.Routes.Get("/users/profile", app.Handlers.CurrentUserProfile)
	app.App.Routes.Get("/users/forgot-password", app.Handlers.ForgotPassword)
	app.App.Routes.Post("/users/forgot-password", app.Handlers.PostForgotPassword)
	app.App.Routes.Get("/users/reset-password", app.Handlers.ResetPasswordForm)
	app.App.Routes.Post("/users/reset-password", app.Handlers.PostResetPassword)

	// API routes
	app.App.Routes.Mount("/api", app.ApiRoutes())

	// Static file server
	fileServer := http.FileServer(http.Dir("./public"))
	app.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return app.App.Routes

}
