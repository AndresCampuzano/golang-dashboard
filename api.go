package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	router.HandleFunc("/user", makeHTTPHandlerFunc(s.handleUser))
	router.HandleFunc("/users", withJWTAuth(makeHTTPHandlerFunc(s.handleUsers), s.store))
	//router.HandleFunc("/users/{id}", withJWTAuth(makeHTTPHandlerFunc(s.handleUsersAndID), s.store))

	log.Println("JSON API server running on port: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
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
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleUser handles user creation.
func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
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

// handleGetUsers handles requests to retrieve all users.
func (s *APIServer) handleGetUsers(w http.ResponseWriter, _ *http.Request) error {
	users, err := s.store.GetUsers()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, users)
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

// WriteJSON writes the given data as JSON to the HTTP response with the provided status code.
// It sets the "Content-Type" header to "application/json; charset=utf-8".
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

// apiFunc is a type representing a function that handles HTTP requests and returns an error.
type apiFunc func(http.ResponseWriter, *http.Request) error

// ApiError represents an error response in JSON format.
type ApiError struct {
	Error string `json:"error"`
}

// makeHTTPHandlerFunc creates an HTTP handler function from the given apiFunc.
// It calls the provided function f to handle HTTP requests, and if an error occurs, it writes
// the error response as JSON with status code http.StatusBadRequest.
func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			err := WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}
}

// getID extracts the ID parameter from the URL path of the HTTP request r.
// It returns the extracted ID and an error if the ID is invalid or not found in the request.
func getID(r *http.Request) (string, error) {
	id := mux.Vars(r)["id"]

	_, err := uuid.Parse(id)
	if err != nil {
		return id, fmt.Errorf("invalid user id %s: %v", id, err)
	}
	return id, nil
}
