package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
)

// createJWT generates a JSON Web Token (JWT) containing the specified user information.
// It returns the signed JWT token as a string and any error encountered during token generation.
func createJWT(user *User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"email":     user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

// permissionDeniedError
func permissionDeniedError(w http.ResponseWriter) {
	err := WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
	if err != nil {
		log.Fatal(err)
		return
	}
}

// withJWTAuth adds JWT authentication to the provided HTTP handler.
// It validates the included JWT and authorizes the request.
// If the JWT is invalid or the request is unauthorized, it responds with a permission denied error.
// Returns an HTTP handler that wraps the original handler.
func withJWTAuth(fn http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		tokenString := r.Header.Get("Authorization")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDeniedError(w)
			return
		}

		if !token.Valid {
			permissionDeniedError(w)
			return
		}

		if err != nil {
			err := WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "invalid token"})
			if err != nil {
				log.Fatal(err)
				return
			}
			return
		}

		fn(w, r)
	}
}

// validateJWT validates the given JWT token string. It verifies the signature
// and checks if the token is well-formed and valid.
func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
