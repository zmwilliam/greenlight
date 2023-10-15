package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zmwilliam/greenlight/internal/data"
)

func (app application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
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

	if err := app.writeJSON(w, http.StatusOK, envelope{"movie": movie}); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
