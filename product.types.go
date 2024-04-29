package main

import "time"

type Product struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Price           int       `json:"price"`
	Image           string    `json:"image"`
	AvailableColors []string  `json:"available_colors"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name            string   `json:"name"`
	Price           int      `json:"price"`
	Image           string   `json:"image"`
	AvailableColors []string `json:"available_colors"`
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
