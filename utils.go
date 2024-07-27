package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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

func colorFromLocalConstants(color string) struct {
	bgColor   string
	textColor string
} {
	// Implement color lookup logic here
	// Return background color and text color based on the provided color
	return struct {
		bgColor   string
		textColor string
	}{"#ffffff", "#000000"}
}

func currencyFormat(price string) string {
	// Implement currency formatting logic here
	return price
}

func findProductVariation(products []*Product, pID string) *Product {
	for _, x := range products {
		if x.ID == pID {
			return x
		}
	}
	return nil
}

// formatCurrency formats an integer value as a currency string.
// Examples:
//
//	formattedValue := formatCurrency(45000, "COP")
//	// Output: COP 45,000
//	formattedValue := formatCurrency(120000, "COP")
//	// Output: COP 120,000
//	formattedValue := formatCurrency(60000, "COP")
//	// Output: COP 60,000
func formatCurrency(value int, currency string) string {
	// Convert int to string
	valueStr := strconv.Itoa(value)

	// Split value string into integer and decimal parts
	parts := strings.Split(valueStr, ".")

	// Format integer part with commas for thousands
	integerPart := parts[0]
	var formattedInteger string
	for i := len(integerPart); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		formattedInteger = integerPart[start:i] + formattedInteger
		if start != 0 {
			formattedInteger = "," + formattedInteger
		}
	}

	// Construct formatted currency string
	formattedValue := currency + " " + formattedInteger
	if len(parts) > 1 {
		formattedValue += "." + parts[1]
	}

	return formattedValue
}

// Color represents a color option
type Color struct {
	Label     string
	Color     string
	TextColor string
}

// Colors contains available color options
var Colors = []Color{
	{
		Label: "Rojo",
		Color: "#a42222",
	},
	{
		Label:     "Azul",
		Color:     "#0cadde",
		TextColor: "#070707",
	},
	{
		Label: "Verde",
		Color: "#1b8c1b",
	},
	{
		Label:     "Amarillo",
		Color:     "#e0e010",
		TextColor: "#070707",
	},
	{
		Label: "Morado",
		Color: "#b622b6",
	},
	{
		Label:     "Blanco",
		Color:     "#ffffff",
		TextColor: "#000000",
	},
	{
		Label: "Cafe",
		Color: "#52302a",
	},
	{
		Label:     "Naranja",
		Color:     "#e07c10",
		TextColor: "#070707",
	},
	{
		Label: "Rosado",
		Color: "#d272d5",
	},
	{
		Label:     "Negro",
		Color:     "#070707",
		TextColor: "#fafafa",
	},
	{
		Label:     "Otro",
		Color:     "#a6a6a6",
		TextColor: "#232323",
	},
	{
		Label:     "Personalizado",
		Color:     "#000000",
		TextColor: "#ffffff",
	},
}

// ColorFromLocalConstants retrieves the color from the constants based on the label
func ColorFromLocalConstants(label string) (string, string) {
	for _, color := range Colors {
		if color.Label == label {
			return color.Color, color.TextColor
		}
	}
	// Default values
	return "#000000", "#ffffff"
}
