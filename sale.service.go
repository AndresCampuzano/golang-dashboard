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

	if err := server.store.CreateSale(sale); err != nil {
		return err
	}

	// Recovering product from DB
	//createdProduct, err := server.store.GetProductByID(product.ID)
	//if err != nil {
	//	return err
	//}

	// Return the newly created product in the response
	return WriteJSON(w, http.StatusOK, nil)
}
