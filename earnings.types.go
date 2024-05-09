package main

import "time"

type ExpensesSummary struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

type CitiesSummary struct {
	Name  string `json:"name"`
	Sales int    `json:"sales"`
}

type DepartmentsSummary struct {
	Name  string `json:"name"`
	Sales int    `json:"sales"`
}

type PurchasedProductsSummary struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Color    string `json:"color"`
	Price    int    `json:"price"`
	Image    string `json:"image"`
	Quantity int    `json:"quantity"`
}

type Earnings struct {
	SortByMonth        time.Time         `json:"sort_by_month"`
	ExpensesSummary    []ExpensesSummary `json:"expenses_summary"`
	AllExpensesInMonth []struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Type        string  `json:"type"`
		Description string  `json:"description"`
		Currency    string  `json:"currency"`
		CreatedAt   string  `json:"created_at"` // FIXME: fix this
		UpdatedAt   string  `json:"updated_at"` // FIXME: fix this
	} `json:"all_expenses_in_month"`
	Income                        float64                    `json:"income"`
	CopExpense                    float64                    `json:"cop_expense"`
	Earnings                      float64                    `json:"earnings"`
	TotalSalesInMonth             int                        `json:"total_sales_in_month"`
	TotalProductVariationsInMonth int                        `json:"total_product_variations_in_month"`
	Cities                        []CitiesSummary            `json:"cities"`
	Departments                   []DepartmentsSummary       `json:"departments"`
	PurchasedProducts             []PurchasedProductsSummary `json:"purchased_products"`
}
