package main

import (
	"net/http"
	"time"

	"github.com/zmwilliam/greenlight/internal/data"
)

func (app application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.writeError(w, http.StatusNotFound)
		return
	}

	movie := data.Movie{
		ID:        id,
		Title:     "The Shawnshank Redemption",
		Runtime:   142,
		Genres:    []string{"drama", "prision", "friendship"},
		Year:      1994,
		CreatedAt: time.Now(),
		Version:   1,
	}

	if err := app.writeJSON(w, envelope{"movie": movie}); err != nil {
		app.writeError(w, http.StatusInternalServerError)
	}
}
