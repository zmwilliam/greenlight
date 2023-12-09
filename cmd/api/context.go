package main

import (
	"context"
	"net/http"

	"github.com/zmwilliam/greenlight/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (*application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (*application) contextGetUser(r *http.Request) *data.User {
	if user, ok := r.Context().Value(userContextKey).(*data.User); ok {
		return user
	}
	panic("missing user value in request context")
}
