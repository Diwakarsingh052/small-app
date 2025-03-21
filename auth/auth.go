package auth

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int

const Key ctxKey = 1

type Auth struct {
	publicKey *rsa.PublicKey
}

// New function would return the implementation of the auth struct
// it is common approach in go community to initialize a struct with private fields
func New(publicKey *rsa.PublicKey) (*Auth, error) {
	if publicKey == nil {
		return nil, errors.New("public key is required")
	}
	return &Auth{publicKey: publicKey}, nil
}

func (a *Auth) ValidateToken(token string) (jwt.RegisteredClaims, error) {
	//read the public every time when this function is called
	var c jwt.RegisteredClaims
	tkn, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}
	if !tkn.Valid {
		return jwt.RegisteredClaims{}, err
	}
	return c, nil
}
