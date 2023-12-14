package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zmwilliam/greenlight/internal/data"
	"github.com/zmwilliam/greenlight/internal/validator"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

type QueryParams struct {
	params url.Values
}

func (q QueryParams) GetString(key, defaultValue string) string {
	if s := q.params.Get(key); s != "" {
		return s
	}
	return defaultValue
}

func (q QueryParams) GetCSV(key string, defaultValue []string) []string {
	if s := q.params.Get(key); s != "" {
		return strings.Split(s, ",")
	}

	return defaultValue
}

func (q QueryParams) GetInt(key string, defaultValue int) (int, error) {
	s := q.params.Get(key)
	if s == "" {
		return defaultValue, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue, err
	}

	return i, nil
}

func NewQueryParams(r *http.Request) QueryParams {
	return QueryParams{params: r.URL.Query()}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	var err error
	v := validator.New()
	qs := NewQueryParams(r)

	input.Title = qs.GetString("title", "")
	input.Genres = qs.GetCSV("genres", []string{})
	if input.Filters.Page, err = qs.GetInt("page", defaultPageNum); err != nil {
		v.AddError("page", "invalid query param, must be integer")
	}
	if input.Filters.PageSize, err = qs.GetInt("page_size", defaultPageSize); err != nil {
		v.AddError("page_size", "invalid query param, must be integer")
	}
	input.Filters.Sort = qs.GetString("sort", "id")
	input.Filters.SortSafelist = []string{
		"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime",
	}

	if input.Filters.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, meta, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": meta}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
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
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/movies/%d", movie.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	updateMovie := func(movie *data.Movie) {
		movie.Title = input.Title
		movie.Year = input.Year
		movie.Runtime = input.Runtime
		movie.Genres = input.Genres
	}

	app.readValidateAndUpdateMovie(w, r, &input, updateMovie)
}

func (app application) patchMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	patchUpdate := func(movie *data.Movie) {
		if input.Title != nil {
			movie.Title = *input.Title
		}
		if input.Year != nil {
			movie.Year = *input.Year
		}
		if input.Runtime != nil {
			movie.Runtime = *input.Runtime
		}
		if input.Genres != nil {
			movie.Genres = input.Genres
		}
	}

	app.readValidateAndUpdateMovie(w, r, &input, patchUpdate)
}

func (app application) readValidateAndUpdateMovie(
	w http.ResponseWriter,
	r *http.Request,
	readInto any,
	update_attrs func(m *data.Movie),
) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.readJSON(w, r, &readInto)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	update_attrs(movie)

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if err = app.models.Movies.Delete(id); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(
		w,
		http.StatusCreated,
		envelope{"message": "movie succesfully deleted"},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
