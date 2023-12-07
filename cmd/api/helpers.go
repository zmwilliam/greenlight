package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type envelope map[string]interface{}

func (*application) readIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if id, err := strconv.ParseInt(idParam, 10, 64); err == nil {
		return id, nil
	}
	return 0, errors.New("invalid id parameter")
}

func (*application) writeJSON(
	w http.ResponseWriter,
	status int,
	data envelope,
	headers http.Header,
) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, values := range headers {
		w.Header()[key] = values
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(append(js, '\n'))
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	max_bytes := 1_048_567
	r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("Body contains badly-formed JSON at char %d", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf(
					"body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field,
				)
			}
			return fmt.Errorf(
				"body contains incorrect JSON type at char %d",
				unmarshalTypeError.Offset,
			)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", max_bytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must only contain a single JSON value")
	}

	return nil
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}
