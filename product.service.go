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

func (server *APIServer) handleGetProducts(w http.ResponseWriter, _ *http.Request) error {
	products, err := server.store.GetProducts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, products)
}

func (server *APIServer) handleGetProductByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	product, err := server.store.GetProductByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, product)
}

func (server *APIServer) handleUpdateProduct(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	oldProduct, err := server.store.GetProductByID(id)
	if err != nil {
		return err
	}

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return err
	}

	// Check if the image has changed, if it changed, it will always be a new base64 file
	if oldProduct.Image != product.Image {
		// delete the old aws image
		err = BucketBasics.DeleteFile(BucketBasics{S3Client: server.s3Client}, oldProduct.Image)
		if err != nil {
			return err
		}
		// Upload the new base 64 image
		imageUrl, err := BucketBasics.UploadFile(BucketBasics{S3Client: server.s3Client}, product.Image)
		if err != nil {
			return err
		}

		product.Image = imageUrl

	}

	product.ID = id

	if err := server.store.UpdateProduct(&product); err != nil {
		return err
	}

	// Retrieve the updated information from the database to get the most up-to-date data
	updatedProduct, err := server.store.GetProductByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, updatedProduct)
}

func (server *APIServer) handleDeleteProduct(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	product, err := server.store.GetProductByID(id)
	if err != nil {
		return err
	}
	// delete the old aws image
	err = BucketBasics.DeleteFile(BucketBasics{S3Client: server.s3Client}, product.Image)
	if err != nil {
		return err
	}

	if err := server.store.DeleteProduct(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"deleted": id})
}
