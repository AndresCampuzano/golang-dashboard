package main

import "time"

type Sale struct {
	ID                       string    `json:"id"`
	CustomerID               string    `json:"customer_id"`
	CustomerName             string    `json:"customer_name"`
	CustomerInstagramAccount string    `json:"customer_instagram_account"`
	CustomerPhone            int       `json:"customer_phone"`
	CustomerAddress          string    `json:"customer_address"`
	CustomerCity             string    `json:"customer_city"`
	CustomerDepartment       string    `json:"customer_department"`
	CustomerComments         string    `json:"customer_comments"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type SaleCustomerRequest struct {
	CustomerID               string `json:"customer_id"`
	CustomerName             string `json:"customer_name"`
	CustomerInstagramAccount string `json:"customer_instagram_account"`
	CustomerPhone            int    `json:"customer_phone"`
	CustomerAddress          string `json:"customer_address"`
	CustomerCity             string `json:"customer_city"`
	CustomerDepartment       string `json:"customer_department"`
	CustomerComments         string `json:"customer_comments"`
}

type ProductVariations struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Color     string    `json:"color"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductVariationsRequest struct {
	ProductID string `json:"product_id"`
	Color     string `json:"color"`
	Price     int    `json:"price"`
}

type SaleProducts struct {
	SaleID             string `json:"sale_id"`
	ProductVariationID string `json:"product_variation_id"`
}
