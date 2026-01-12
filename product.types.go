package main

import "time"

type CatalogVariant struct {
	ID        string    `json:"id"`
	ColorHex  string    `json:"color_hex"`
	ColorName string    `json:"color_name"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Price           int              `json:"price"`
	Image           string           `json:"image"`
	AvailableColors []string         `json:"available_colors"`
	Description     *string          `json:"description,omitempty"`
	IsCatalogReady  bool             `json:"is_catalog_ready"`
	CatalogVariants []CatalogVariant `json:"catalog_variants,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type CreateProductRequest struct {
	Name            string           `json:"name"`
	Price           int              `json:"price"`
	Image           string           `json:"image"`
	AvailableColors []string         `json:"available_colors"`
	Description     *string          `json:"description,omitempty"`
	IsCatalogReady  *bool            `json:"is_catalog_ready,omitempty"`
	CatalogVariants []CatalogVariant `json:"catalog_variants,omitempty"`
}

func NewProduct(
	name string,
	price int,
	image string,
	availableColors []string,
) (*Product, error) {
	return &Product{
		Name:            name,
		Price:           price,
		Image:           image,
		AvailableColors: availableColors,
	}, nil
}
