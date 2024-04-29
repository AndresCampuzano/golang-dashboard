package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
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

// ConvertToDBArray converts a slice of strings to a format suitable for PostgreSQL array type.
// Example:
//
//	input: []string{"foo", "bar"}
//	output: "{\"foo\",\"bar\"}"
func ConvertToDBArray(localString []string) string {
	var quoted []string
	for _, item := range localString {
		escapedItem := strings.ReplaceAll(item, `"`, `\"`)
		quoted = append(quoted, `"`+escapedItem+`"`)
	}
	return "{" + strings.Join(quoted, ",") + "}"
}

// ConvertFromDBArray converts a PostgreSQL array string to a slice of strings.
// Example:
//
//	input: "{\"foo\",\"bar\"}"
//	output: []string{"foo", "bar"}
func ConvertFromDBArray(dbArray string) []string {
	// Trim the curly braces from the input string
	dbArray = strings.Trim(dbArray, "{}")
	if dbArray == "" {
		return nil
	}
	// Split the input string by commas to extract individual elements
	return strings.Split(dbArray, ",")
}
