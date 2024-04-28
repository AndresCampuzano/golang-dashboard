package main

import (
	"encoding/json"
	"net/http"
)

func (s *APIServer) handleCreateCustomer(w http.ResponseWriter, r *http.Request) error {
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

	if err := s.store.CreateCustomer(customer); err != nil {
		return err
	}

	// Recovering customer from DB
	createdCustomer, err := s.store.GetCustomerByID(customer.ID)
	if err != nil {
		return err
	}

	// Return the newly created customer in the response
	return WriteJSON(w, http.StatusOK, createdCustomer)
}
