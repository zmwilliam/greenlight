package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (application) notImpletedYetHandler(handlerName string) http.HandlerFunc {
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

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		r.Route("/movies", func(r chi.Router) {
			r.Get("/", app.notImpletedYetHandler("listMovies"))
			r.Post("/", app.notImpletedYetHandler("createMovie"))

			r.Get("/{id}", app.notImpletedYetHandler("showMovie"))
			r.Put("/{id}", app.notImpletedYetHandler("editMovie"))
			r.Delete("/{id}", app.notImpletedYetHandler("deleteMovie"))
		})
	})

	return r
}
