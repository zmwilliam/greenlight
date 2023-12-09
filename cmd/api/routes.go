package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		r.Route("/movies", func(r chi.Router) {
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

	// adding recover panic middleware,
	// this can also be done using chi "Recoverer middleware"
	// https://go-chi.io/#/pages/middleware?id=recoverer
	return app.recoverPanic(app.rateLimit(r))
}
