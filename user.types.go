package main

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type User struct {
	ID                string    `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"-"`
	CreatedAt         time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (a *User) ValidatePassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw))
	if err != nil {
		return false
	}

	return true
}

func NewUser(
	firstName string,
	lastName string,
	email string,
	password string,
) (*User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         firstName,
		LastName:          lastName,
		Email:             email,
		EncryptedPassword: string(encryptedPassword),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
