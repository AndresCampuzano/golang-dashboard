package main

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *PostgresStore) CreateCustomersTable() error {
	// Create the table if it doesn't exist
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS customers (
            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            instagram_account VARCHAR(255) NOT NULL,
            phone BIGINT,
            address VARCHAR(255) NOT NULL,
            city VARCHAR(255) NOT NULL,
            department VARCHAR(255) NOT NULL,
            comments VARCHAR(255) NOT NULL,
            cc VARCHAR(255) NOT NULL, 
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return err
	}

	// Check if the trigger already exists
	var triggerExists bool
	err = s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM pg_trigger 
         	WHERE tgname = 'customers_updated_at_trigger' 
           	AND tgrelid = 'customers'::regclass)
           `).Scan(&triggerExists)
	if err != nil {
		return err
	}

	// Create the trigger only if it doesn't exist
	if !triggerExists {
		// Create the trigger for updated_at
		_, err = s.db.Exec(`
            CREATE OR REPLACE FUNCTION update_timestamp()
            RETURNS TRIGGER AS $$
            BEGIN
                NEW.updated_at = NOW();
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;
            
            CREATE TRIGGER customers_updated_at_trigger
            BEFORE UPDATE ON customers
            FOR EACH ROW
            EXECUTE FUNCTION update_timestamp();
        `)
		if err != nil {
			return err
		}
	}

	// Check if the createdAt trigger already exists
	err = s.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM pg_trigger 
            WHERE tgname = 'customers_created_at_trigger' 
            AND tgrelid = 'customers'::regclass)
    `).Scan(&triggerExists)
	if err != nil {
		return err
	}

	// Create the trigger only if it doesn't exist
	if !triggerExists {
		// Create the trigger for created_at
		_, err = s.db.Exec(`
            CREATE OR REPLACE FUNCTION set_created_at()
            RETURNS TRIGGER AS $$
            BEGIN
                NEW.created_at = NOW();
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;
            
            CREATE TRIGGER customers_created_at_trigger
            BEFORE INSERT ON customers
            FOR EACH ROW
            EXECUTE FUNCTION set_created_at();
        `)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStore) CreateCustomer(customer *Customer) error {
	query := `
        INSERT INTO customers (
			name, 
			instagram_account, 
			phone, 
			address, 
			city, 
			department, 
			comments,
			cc,
			created_at, 
			updated_at 
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
        RETURNING id
    `

	var id string
	err := s.db.QueryRow(
		query,
		customer.Name,
		customer.InstagramAccount,
		customer.Phone,
		customer.Address,
		customer.City,
		customer.Department,
		customer.Comments,
		customer.Cc,
		customer.CreatedAt,
		customer.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}

	// Set the ID of the inserted customer
	customer.ID = id

	return nil
}

func (s *PostgresStore) GetCustomerByID(id string) (*Customer, error) {
	rows, err := s.db.Query("SELECT * FROM customers WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	for rows.Next() {
		return scanIntoCustomers(rows)
	}

	return nil, fmt.Errorf("customer [%s] not found", id)
}

func scanIntoCustomers(rows *sql.Rows) (*Customer, error) {
	customer := new(Customer)
	err := rows.Scan(
		&customer.ID,
		&customer.Name,
		&customer.InstagramAccount,
		&customer.Phone,
		&customer.Address,
		&customer.City,
		&customer.Department,
		&customer.Comments,
		&customer.Cc,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	return customer, err
}

func (s *PostgresStore) GetCustomers() ([]*Customer, error) {
	rows, err := s.db.Query("SELECT * FROM customers")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var customers []*Customer
	for rows.Next() {
		customer, err := scanIntoCustomers(rows)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func (s *PostgresStore) GetCustomersLast3Months() ([]*Customer, error) {
	rows, err := s.db.Query(`
		SELECT * FROM customers WHERE created_at >= (NOW() - INTERVAL '3 months') ORDER BY created_at DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}(rows)
	var customers []*Customer
	for rows.Next() {
		customer, err := scanIntoCustomers(rows)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (s *PostgresStore) UpdateCustomer(customer *Customer) error {
	query := `
		UPDATE customers
		SET 
		    name = $1, 
		    instagram_account = $2, 
		    phone = $3, 
		    address = $4, 
		    city = $5, 
		    department = $6, 
		    comments = $7,
		    cc = $8
		WHERE id = $9
	`

	_, err := s.db.Exec(
		query,
		customer.Name,
		customer.InstagramAccount,
		customer.Phone,
		customer.Address,
		customer.City,
		customer.Department,
		customer.Comments,
		customer.Cc,
		customer.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteCustomer(id string) error {
	_, err := s.db.Exec("DELETE FROM customers WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
