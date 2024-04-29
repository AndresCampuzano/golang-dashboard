package main

import (
	"database/sql"
	"fmt"
)

func (s *PostgresStore) CreateProductsTable() error {
	// Create the table if it doesn't exist
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            price BIGINT NOT NULL,
            image VARCHAR(255) NOT NULL,
            available_colors VARCHAR(20)[] NOT NULL DEFAULT '{}'::VARCHAR(20)[],
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
         	WHERE tgname = 'products_updated_at_trigger' 
           	AND tgrelid = 'products'::regclass)
           `).Scan(&triggerExists)
	if err != nil {
		return err
	}

	// Create the trigger only if it doesn't exist
	if !triggerExists {
		// Create the trigger
		_, err = s.db.Exec(`
            CREATE OR REPLACE FUNCTION update_timestamp()
            RETURNS TRIGGER AS $$
            BEGIN
                NEW.updated_at = NOW();
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;
            
            CREATE TRIGGER products_updated_at_trigger
            BEFORE UPDATE ON products
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
            WHERE tgname = 'products_created_at_trigger' 
            AND tgrelid = 'products'::regclass)
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
            
            CREATE TRIGGER products_created_at_trigger
            BEFORE INSERT ON products
            FOR EACH ROW
            EXECUTE FUNCTION set_created_at();
        `)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStore) CreateProduct(product *Product) error {

	query := `
        INSERT INTO products (
            name,
            price,
            image,
            available_colors,
            created_at,
            updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	availableColorsDB := ConvertToDBArray(product.AvailableColors)
	var id string
	err := s.db.QueryRow(
		query,
		product.Name,
		product.Price,
		product.Image,
		availableColorsDB,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}

	// Set the ID of the inserted product
	product.ID = id

	return nil
}

func (s *PostgresStore) GetProductByID(id string) (*Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoProducts(rows)
	}

	return nil, fmt.Errorf("product [%s] not found", id)
}

func scanIntoProducts(rows *sql.Rows) (*Product, error) {
	product := new(Product)
	var availableColorsDB string
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Image,
		&availableColorsDB,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	product.AvailableColors = ConvertFromDBArray(availableColorsDB)

	return product, nil
}
