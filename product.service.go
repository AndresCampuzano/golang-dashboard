package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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

	// Set new optional fields
	product.Description = req.Description
	if req.IsCatalogReady != nil {
		product.IsCatalogReady = *req.IsCatalogReady
	}

	// Process catalog variants if provided
	if len(req.CatalogVariants) > 0 {
		processedVariants, err := server.processCatalogVariants(req.CatalogVariants)
		if err != nil {
			return err
		}
		product.CatalogVariants = processedVariants
	}

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

	// Process catalog variants if provided
	if len(product.CatalogVariants) > 0 {
		processedVariants, err := server.processCatalogVariants(product.CatalogVariants)
		if err != nil {
			return err
		}
		product.CatalogVariants = processedVariants
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

func (server *APIServer) processCatalogVariants(variants []CatalogVariant) ([]CatalogVariant, error) {
	now := time.Now()
	for i := range variants {
		// Generate UUID if empty
		if variants[i].ID == "" {
			variants[i].ID = uuid.New().String()
		}
		// Set timestamps
		if variants[i].CreatedAt.IsZero() {
			variants[i].CreatedAt = now
		}
		variants[i].UpdatedAt = now

		// Upload image to S3 if base64
		if strings.HasPrefix(variants[i].Image, "data:image") {
			url, err := BucketBasics.UploadFile(BucketBasics{S3Client: server.s3Client}, variants[i].Image)
			if err != nil {
				return nil, err
			}
			variants[i].Image = url
		}
	}
	return variants, nil
}

func (server *APIServer) handleGetPublicProducts(w http.ResponseWriter, _ *http.Request) error {
	products, err := server.store.GetCatalogProducts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, products)
}

func (server *APIServer) handlePublicProducts(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return server.handleGetPublicProducts(w, r)
	}
	return WriteJSON(w, http.StatusMethodNotAllowed, apiError{Error: "unsupported method: " + r.Method})
}
