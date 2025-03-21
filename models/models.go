package models

// User struct would be used to store values inside a database or store fetched values from db
type User struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"` // giving the name of the field in the json output
	Age          int64  `json:"age"`
	PasswordHash string `json:"password_hash"` //ignore the field while creating the json
}

type NewUser struct {
	Name     string `json:"name" validate:"required,max=60,min=3"`
	Email    string `json:"email" validate:"required,email,max=60,min=5"`
	Age      int64  `json:"age" validate:"required,min=18"`
	Password string `json:"password" validate:"required,max=60"`
}
