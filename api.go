package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"strings"
)

// APIServer represents an HTTP server for handling API requests.
type APIServer struct {
	listenAddr string
	store      Storage
	s3Client   *s3.Client
	Router     *mux.Router
}

// NewAPIServer creates a new instance of APIServer.
func NewAPIServer(listenAddr string, store Storage, s3 *s3.Client) *APIServer {
	router := mux.NewRouter()

	server := &APIServer{
		listenAddr: listenAddr,
		store:      store,
		s3Client:   s3,
		Router:     router,
	}

	router.HandleFunc("/api/healthcheck", makeHTTPHandlerFunc(server.handleHealth))
	router.HandleFunc("/api/login", makeHTTPHandlerFunc(server.handleLogin))
	//router.HandleFunc("/api/signup", makeHTTPHandlerFunc(server.HandleSignUp))
	router.HandleFunc("/api/users", withJWTAuth(makeHTTPHandlerFunc(server.handleUsers), server.store))
	router.HandleFunc("/api/users/{id}", withJWTAuth(makeHTTPHandlerFunc(server.handleUsersWithID), server.store))
	router.HandleFunc("/api/customers", withJWTAuth(makeHTTPHandlerFunc(server.handleCustomers), server.store))
	router.HandleFunc("/api/customers/{id}", withJWTAuth(makeHTTPHandlerFunc(server.handleCustomersWithID), server.store))
	router.HandleFunc("/api/products", withJWTAuth(makeHTTPHandlerFunc(server.handleProducts), server.store))
	router.HandleFunc("/api/products/{id}", withJWTAuth(makeHTTPHandlerFunc(server.handleProductsWithID), server.store))
	router.HandleFunc("/api/sales", withJWTAuth(makeHTTPHandlerFunc(server.handleSales), server.store))
	router.HandleFunc("/api/sales/{id}", withJWTAuth(makeHTTPHandlerFunc(server.handleSalesWithID), server.store))
	router.HandleFunc("/api/expenses", withJWTAuth(makeHTTPHandlerFunc(server.handleExpenses), server.store))
	router.HandleFunc("/api/expenses/{id}", withJWTAuth(makeHTTPHandlerFunc(server.handleExpensesWithID), server.store))
	router.HandleFunc("/api/earnings", withJWTAuth(makeHTTPHandlerFunc(server.handleEarnings), server.store))

	return server
}

// Run starts the API server and listens for incoming requests.
func (server *APIServer) Run() {
	log.Println("JSON API server running on port: ", server.listenAddr)

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	origins := strings.Split(allowedOrigins, ",")

	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		Debug:            true,
	})

	handler := c.Handler(server.Router)

	// Use the CORS-wrapped handler as your HTTP server's handler
	err := http.ListenAndServe(server.listenAddr, handler)
	if err != nil {
		log.Fatal(err)
	}
}

// handleHealth sends a 200 status code.
func (server *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleHealthCheck(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
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
	case http.MethodGet:
		return server.handleGetSales(w, r)
	case http.MethodPost:
		return server.handleCreateSale(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleSalesWithID handles get, update and delete requests
func (server *APIServer) handleSalesWithID(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleGetSaleByID(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleExpenses handles get and post requests
func (server *APIServer) handleExpenses(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleGetExpenses(w, r)
	case http.MethodPost:
		return server.handleCreateExpense(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleExpensesWithID handles get, update and delete requests
func (server *APIServer) handleExpensesWithID(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleGetExpenseByID(w, r)
	case http.MethodPut:
		return server.handleUpdateExpense(w, r)
	case http.MethodDelete:
		return server.handleDeleteExpense(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleEarnings handles get by month
func (server *APIServer) handleEarnings(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return server.handleGetEarnings(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}
