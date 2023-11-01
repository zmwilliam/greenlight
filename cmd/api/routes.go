package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app application) routes() (router http.Handler) {
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
	})

	// adding recover panic middleware,
	// this can also be done using chi "Recoverer middleware"
	// https://go-chi.io/#/pages/middleware?id=recoverer
	router = app.recoverPanic(r)
	return
}
