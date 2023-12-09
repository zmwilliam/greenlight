package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.recoverPanic)
	r.Use(app.rateLimit)
	r.Use(app.authenticate)

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		r.Route("/movies", func(r chi.Router) {
			r.Use(app.requireActivatedUser)
			r.Get("/", app.listMoviesHandler)
			r.Post("/", app.createMovieHandler)

			r.Get("/{id}", app.showMovieHandler)
			r.Put("/{id}", app.updateMovieHandler)
			r.Patch("/{id}", app.patchMovieHandler)
			r.Delete("/{id}", app.deleteMovieHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.registerUserHandler)
			r.Put("/activated", app.activateUserHandler)
		})

		r.Post("/tokens/authentication", app.createAuthTokenHandler)
	})

	return r
}
