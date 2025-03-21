package middleware

import (
	"errors"
	"rest-api/auth"
)

// when we need to inject dependency from another package
// we would create a struct that would contain the required dependency

type Mid struct {
	a *auth.Auth
}

func NewMid(a *auth.Auth) (*Mid, error) {
	if a == nil {
		return nil, errors.New("auth is nil")
	}
	return &Mid{a: a}, nil
}
