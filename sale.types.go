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
	Image string `json:"image"`
	Name  string `json:"name"`
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
	CustomerTotalPurchases   int                         `json:"customer_total_purchases"`
	CreatedAt                string                      `json:"created_at"`
	UpdatedAt                string                      `json:"updated_at"`
	ProductVariations        []ProductVariationsResponse `json:"product_variations"`
	OtherSales               []SaleResponse              `json:"other_sales"`
	// OtherSales: it returns null inside nested data,
	// so it could be an empty slice, with data or 'null' after parsing it
}

type SaleResponseSortedByMonth struct {
	SaleResponse
	SortByMonth time.Time `json:"sort_by_month"`
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
