package main

import (
	"encoding/json"
	"net/http"
)

func (server *APIServer) handleCreateCustomer(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateCustomerRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	customer, err := NewCustomer(
		req.Name,
		req.InstagramAccount,
		req.Phone,
		req.Address,
		req.City,
		req.Department,
		req.Comments,
	)
	if err != nil {
		return err
	}

	if err := server.store.CreateCustomer(customer); err != nil {
		return err
	}

	// Recovering customer from DB
	createdCustomer, err := server.store.GetCustomerByID(customer.ID)
	if err != nil {
		return err
	}

	// Return the newly created customer in the response
	return WriteJSON(w, http.StatusOK, createdCustomer)
}

func (server *APIServer) handleGetCustomers(w http.ResponseWriter, _ *http.Request) error {
	users, err := server.store.GetCustomers()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, users)
}

//func (server *APIServer) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) error {
//	id, err := getID(r)
//	if err != nil {
//		return err
//	}
//
//	customer, err := server.store.GetCustomerByID(id)
//	if err != nil {
//		return err
//	}
//
//	return WriteJSON(w, http.StatusOK, customer)
//}
//
//
