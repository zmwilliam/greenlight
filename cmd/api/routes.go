package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (application) notImplementedYetHandler(handlerName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		if id := chi.URLParam(r, "id"); id != "" {
			fmt.Fprintf(w, "%s %s /%s not implemented yet!\n", handlerName, r.Method, id)
		} else {
			fmt.Fprintf(w, "%s not implemented yet!\n", handlerName)
		}
	}
}

func (app application) routes() chi.Router {
	r := chi.NewRouter()

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		r.Route("/movies", func(r chi.Router) {
			r.Get("/", app.notImplementedYetHandler("listMovies"))
			r.Post("/", app.createMovieHandler)

			r.Get("/{id}", app.showMovieHandler)
			r.Put("/{id}", app.updateMovieHandler)
			r.Delete("/{id}", app.notImplementedYetHandler("deleteMovie"))
		})
	})

	return r
}
