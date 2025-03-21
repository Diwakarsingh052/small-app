package models

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CacheStore map[string]User

type Conn struct {
	store CacheStore
}

func NewConn() *Conn {
	return &Conn{store: make(CacheStore, 100)}
}

// ErrNotFound is an error returned when the requested user data is not found.
var ErrNotFound = errors.New("user not found")

func (c *Conn) CreateUser(n NewUser) (User, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	_, ok := c.store[n.Email]
	if ok {
		return User{}, errors.New("user already exists")
	}
	us := User{
		Id:           uuid.NewString(),
		Email:        n.Email,
		Name:         n.Name,
		Age:          n.Age,
		PasswordHash: string(passHash),
	}
	c.store = CacheStore{us.Email: us}
	return us, nil

}

// FetchUser retrieves a User object from the in-memory database by its ID.
// It returns the User object and an error, which is nil if the user is found and ErrDataNotPresent otherwise.
func (c *Conn) FetchUser(userEmail string) (User, error) {
	value, ok := c.store[userEmail]

	// If ok is false, there is no user with the given ID.
	if !ok {
		return User{}, ErrNotFound
	}

	// If ok is true, the user was found and is returned.
	return value, nil
}
