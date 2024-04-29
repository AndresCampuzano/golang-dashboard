package main

import (
	"encoding/json"
	"net/http"
)

func (server *APIServer) handleCreateProduct(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateProductRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	product, err := NewProduct(
		req.Name,
		req.Price,
		req.Image,
		req.AvailableColors,
	)
	if err != nil {
		return err
	}

	imageUrl, err := BucketBasics.UploadFile(BucketBasics{S3Client: server.s3Client}, product.Image)
	if err != nil {
		return err
	}

	product.Image = imageUrl

	if err := server.store.CreateProduct(product); err != nil {
		return err
	}

	// Recovering product from DB
	createdProduct, err := server.store.GetProductByID(product.ID)
	if err != nil {
		return err
	}

	// Return the newly created product in the response
	return WriteJSON(w, http.StatusOK, createdProduct)
}
