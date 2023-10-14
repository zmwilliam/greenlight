package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type envelope map[string]interface{}

func (application) readIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if id, err := strconv.ParseInt(idParam, 10, 64); err == nil {
		return id, nil
	}
	return 0, errors.New("invalid id parameter")
}

func (application) writeJSON(w http.ResponseWriter, status int, data envelope) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(append(js, '\n'))
	return nil
}
