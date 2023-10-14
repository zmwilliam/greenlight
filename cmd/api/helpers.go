package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (application) readIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if id, err := strconv.ParseInt(idParam, 10, 64); err == nil {
		return id, nil
	}
	return 0, errors.New("invalid id parameter")
}

func (application) writeJSON(w http.ResponseWriter, data interface{}) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(append(js, '\n'))
	return nil
}

func (application) writeError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}