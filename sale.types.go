package main

import "time"

type SaleWithProducts struct {
	ID         string              `json:"id"`
	CustomerID string              `json:"customer_id"`
	Products   []ProductVariations `json:"products"`
}

type CreateSaleRequest struct {
	CustomerID string              `json:"customer_id"`
	Products   []ProductVariations `json:"products"`
}

type ProductVariations struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Color     string    `json:"color"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductVariationsResponse struct {
	ID    string `json:"id"`
	Color string `json:"color"`
	Price int    `json:"price"`
}

type SaleResponse struct {
	ID                       string                      `json:"id"`
	CustomerID               string                      `json:"customer_id"`
	CustomerName             string                      `json:"customer_name"`
	CustomerInstagramAccount string                      `json:"customer_instagram_account"`
	CustomerPhone            int                         `json:"customer_phone"`
	CustomerAddress          string                      `json:"customer_address"`
	CustomerCity             string                      `json:"customer_city"`
	CustomerDepartment       string                      `json:"customer_department"`
	CustomerComments         string                      `json:"customer_comments"`
	CustomerCc               string                      `json:"customer_cc"`
	CreatedAt                time.Time                   `json:"created_at"`
	UpdatedAt                time.Time                   `json:"updated_at"`
	ProductVariations        []ProductVariationsResponse `json:"product_variations"`
}

type SaleResponseSortedByMonth struct {
	SaleResponse
	SaleMonth time.Time `json:"sale_moth"`
}

func NewSale(
	customerID string,
	products []ProductVariations,
) (*SaleWithProducts, error) {
	return &SaleWithProducts{
		CustomerID: customerID,
		Products:   products,
	}, nil
}
