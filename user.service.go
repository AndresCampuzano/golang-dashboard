package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (server *APIServer) handleLoginUser(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := server.store.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	if !acc.ValidatePassword(req.Password) {
		return fmt.Errorf("not authorized")
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Email: acc.Email,
		Token: token,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

func (server *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	user, err := NewUser(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return err
	}

	if err := server.store.CreateUser(user); err != nil {
		return err
	}

	tokenString, err := createJWT(user)
	if err != nil {
		return err
	}
	fmt.Println("JWT token: ", tokenString)

	// Recovering user from DB
	createdUser, err := server.store.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	// Return the newly created user in the response
	return WriteJSON(w, http.StatusOK, createdUser)
}

func (server *APIServer) handleGetUsers(w http.ResponseWriter, _ *http.Request) error {
	users, err := server.store.GetUsers()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, users)
}

func (server *APIServer) handleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	account, err := server.store.GetUserByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}
