package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// WriteJSON writes the given data as JSON to the HTTP response with the provided status code.
// It sets the "Content-Type" header to "application/json; charset=utf-8".
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

// apiFunc is a type representing a function that handles HTTP requests and returns an error.
type apiFunc func(http.ResponseWriter, *http.Request) error

// apiError represents an error response in JSON format.
type apiError struct {
	Error string `json:"error"`
}

// makeHTTPHandlerFunc creates an HTTP handler function from the given apiFunc.
// It calls the provided function f to handle HTTP requests, and if an error occurs, it writes
// the error response as JSON with status code http.StatusBadRequest.
func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			err := WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
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
