package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *APIServer) handleLoginUser(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.store.GetUserByEmail(req.Email)
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

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	user, err := NewUser(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateUser(user); err != nil {
		return err
	}

	tokenString, err := createJWT(user)
	if err != nil {
		return err
	}
	fmt.Println("JWT token: ", tokenString)

	// Recovering user from DB
	createdUser, err := s.store.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	// Return the newly created user in the response
	return WriteJSON(w, http.StatusOK, createdUser)
}

// handleGetUsers handles requests to retrieve all users.
func (s *APIServer) handleGetUsers(w http.ResponseWriter, _ *http.Request) error {
	users, err := s.store.GetUsers()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, users)
}

func (s *APIServer) handleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetUserByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}
