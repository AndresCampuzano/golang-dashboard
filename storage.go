package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	// Users
	CreateUser(user *User) error
	GetUsers() ([]*User, error)
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	// Customers
	CreateCustomer(customer *Customer) error
	GetCustomerByID(id string) (*Customer, error)
	GetCustomers() ([]*Customer, error)
	UpdateCustomer(customer *Customer) error
	DeleteCustomer(id string) error
	// Products
	CreateProduct(product *Product) error
	GetProductByID(id string) (*Product, error)
	GetProducts() ([]*Product, error)
	UpdateProduct(product *Product) error
	DeleteProduct(id string) error
	// Sales
	CreateSale(sale *SaleWithProducts) error
	GetSaleByID(id string) (*SaleResponse, error)
	GetSales() ([]*SaleResponse, error)
	GetSalesByMonth() ([]*SaleResponseSortedByMonth, error) // Not in use yet
	// Expenses
	CreateExpense(expense *Expense) error
	GetExpenseByID(id string) (*Expense, error)
	GetExpenses() ([]*Expense, error)
	UpdateExpense(expense *Expense) error
	DeleteExpense(id string) error
	// EarningsSummary
	GetEarnings() ([]*Earnings, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Init() error {
	err := s.CreateUsersTable()
	if err != nil {
		return err
	}

	err = s.CreateCustomersTable()
	if err != nil {
		return err
	}

	err = s.CreateProductsTable()
	if err != nil {
		return err
	}

	err = s.CreateSalesTablesWithRelations()
	if err != nil {
		return err
	}

	err = s.CreateExpensesTable()
	if err != nil {
		return err
	}

	return nil
}
