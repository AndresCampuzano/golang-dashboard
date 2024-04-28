package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(user *User) error
	GetUsers() ([]*User, error)
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateCustomer(customer *Customer) error
	GetCustomerByID(id string) (*Customer, error)
	GetCustomers() ([]*Customer, error)
	UpdateCustomer(customer *Customer) error
	DeleteCustomer(id string) error
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

	return nil
}
