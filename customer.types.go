package main

import "time"

type Customer struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	InstagramAccount string    `json:"instagram_account"`
	Phone            int       `json:"phone"`
	Address          string    `json:"address"`
	City             string    `json:"city"`
	Department       string    `json:"department"`
	Comments         string    `json:"comments"`
	Cc               string    `json:"cc"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateCustomerRequest struct {
	Name             string `json:"name"`
	InstagramAccount string `json:"instagram_account"`
	Phone            int    `json:"phone"`
	Address          string `json:"address"`
	City             string `json:"city"`
	Department       string `json:"department"`
	Comments         string `json:"comments"`
	Cc               string `json:"cc"`
}

//type UpdateCustomerRequest struct {
//	Name             string `json:"name"`
//	InstagramAccount string `json:"instagram_account"`
//	Phone            int    `json:"phone"`
//	Address          string `json:"address"`
//	City             string `json:"city"`
//	Department       string `json:"department"`
//	Comments         string `json:"comments"`
//	Cc         		 string `json:"cc"`
//}

func NewCustomer(
	name string,
	instagramAccount string,
	phone int,
	address string,
	city string,
	department string,
	comments string,
	cc string,
) (*Customer, error) {
	return &Customer{
		Name:             name,
		InstagramAccount: instagramAccount,
		Phone:            phone,
		Address:          address,
		City:             city,
		Department:       department,
		Comments:         comments,
		Cc:               cc,
	}, nil
}
