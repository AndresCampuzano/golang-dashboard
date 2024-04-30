package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// APIServer represents an HTTP server for handling API requests.
type APIServer struct {
	listenAddr string
	store      Storage
	s3Client   *s3.Client
}

// NewAPIServer creates a new instance of APIServer.
func NewAPIServer(listenAddr string, store Storage, s3 *s3.Client) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
		s3Client:   s3,
	}
}

// Run starts the API server and listens for incoming requests.
func (server *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin))
	router.HandleFunc("/signup", makeHTTPHandlerFunc(server.HandleSignUp))
	router.HandleFunc("/users", withJWTAuth(makeHTTPHandlerFunc(server.handleUsers), server.store))
	router.HandleFunc("/users/{id}", withJWTAuth(makeHTTPHandlerFunc(server.handleUsersWithID), server.store))
	router.HandleFunc("/customers", withJWTAuth(makeHTTPHandlerFunc(server.handleCustomers), server.store))
	router.HandleFunc("/customers/{id}", withJWTAuth(makeHTTPHandlerFunc(server.handleCustomersWithID), server.store))
	router.HandleFunc("/products", withJWTAuth(makeHTTPHandlerFunc(server.handleProducts), server.store))
	router.HandleFunc("/products/{id}", withJWTAuth(makeHTTPHandlerFunc(server.handleProductsWithID), server.store))
	router.HandleFunc("/sales", withJWTAuth(makeHTTPHandlerFunc(server.handleSales), server.store))

	log.Println("JSON API server running on port: ", server.listenAddr)

	err := http.ListenAndServe(server.listenAddr, router)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// handleLogin handles user login.
func (server *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		return server.handleLoginUser(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// HandleSignUp handles user sign up.
func (server *APIServer) HandleSignUp(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		return server.handleCreateUser(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleUsers handles user info retrieve.
func (server *APIServer) handleUsers(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleGetUsers(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleUsersWithID handles requests to manage a user by ID.
func (server *APIServer) handleUsersWithID(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleGetUserByID(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleCustomers handles get and post requests.
func (server *APIServer) handleCustomers(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleGetCustomers(w, r)
	case http.MethodPost:
		return server.handleCreateCustomer(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleCustomersWithID handles customer logic.
func (server *APIServer) handleCustomersWithID(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleGetCustomerByID(w, r)
	case http.MethodPut:
		return server.handleUpdateCustomer(w, r)
	case http.MethodDelete:
		return server.handleDeleteCustomer(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleProducts handles get and post requests
func (server *APIServer) handleProducts(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleGetProducts(w, r)
	case http.MethodPost:
		return server.handleCreateProduct(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleCustomersWithID handles get, update and delete requests.
func (server *APIServer) handleProductsWithID(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleGetProductByID(w, r)
	case http.MethodPut:
		return server.handleUpdateProduct(w, r)
	case http.MethodDelete:
		return server.handleDeleteProduct(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleSales handles get and post requests
func (server *APIServer) handleSales(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	//case http.MethodGet:
	//	return server.handleGetSale(w, r)
	case http.MethodPost:
		return server.handleCreateSale(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}
