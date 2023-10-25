package main

import (
	"github.com/go-chi/chi/v5"
)

func (app application) routes() chi.Router {
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

	return r
}
