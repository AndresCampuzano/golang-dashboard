package main

import (
	"encoding/json"
	"net/http"
)

func (server *APIServer) handleCreateSale(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateSaleRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	sale, err := NewSale(
		req.CustomerID,
		req.Products,
	)
	if err != nil {
		return err
	}

	err = server.store.CreateSale(sale)
	if err != nil {
		return err
	}

	// Recovering product from DB
	createdSale, err := server.store.GetSaleByID(sale.ID)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, createdSale)
}

func (server *APIServer) handleGetSales(w http.ResponseWriter, r *http.Request) error {
	sales, err := server.store.GetSales()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, sales)
}

func (server *APIServer) handleGetSaleByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	sale, err := server.store.GetSaleByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, sale)
}
