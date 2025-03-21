package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

/*
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -in private.pem -pubout -out pubkey.pem

go run ./exp/creating-jwt-token.go
*/

func ABC() {
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	claims := jwt.RegisteredClaims{
		Issuer:    "api project",
		Subject:   "101",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(50 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// create encoded token with the claims payload
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	privatePem, err := os.ReadFile("private.pem")
	if err != nil {
		log.Println(err)
		return
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
	if err != nil {
		log.Println(err)
		return
	}
	// singing the token with our private key
	token, err := tkn.SignedString(privateKey)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(token)

}
