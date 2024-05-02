package main

import "time"

type Expense struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateExpenseRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func NewExpense(
	name string,
	price int,
	exType string,
	description string,
) (*Expense, error) {
	return &Expense{
		Name:        name,
		Price:       price,
		Type:        exType,
		Description: description,
	}, nil
}
