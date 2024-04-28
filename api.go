package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// APIServer represents an HTTP server for handling API requests.
type APIServer struct {
	listenAddr string
	store      Storage
}

// NewAPIServer creates a new instance of APIServer.
func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Run starts the API server and listens for incoming requests.
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandlerFunc(s.handleLogin))
	router.HandleFunc("/signup", makeHTTPHandlerFunc(s.HandleSignUp))
	router.HandleFunc("/users", withJWTAuth(makeHTTPHandlerFunc(s.handleUsers), s.store))
	router.HandleFunc("/users/{id}", withJWTAuth(makeHTTPHandlerFunc(s.handleUsersAndID), s.store))

	log.Println("JSON API server running on port: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// handleUser handles user login.
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		return s.handleLoginUser(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// HandleSignUp handles user sign up.
func (s *APIServer) HandleSignUp(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		return s.handleCreateUser(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleUsers handles user info retrieve.
func (s *APIServer) handleUsers(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetUsers(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleUsersAndID handles requests to manage a user by ID.
func (s *APIServer) handleUsersAndID(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetUserByID(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}
