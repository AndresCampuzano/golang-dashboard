package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"
)

func (s *PostgresStore) CreateSalesTablesWithRelations() error {
	// Create the table if it doesn't exist
	_, err := s.db.Exec(`
		-- Create a table to store product variations
		CREATE TABLE IF NOT EXISTS product_variations (
			id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
			product_id UUID REFERENCES products(id) ON DELETE CASCADE,
			color VARCHAR(20) NOT NULL,
			price BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		-- Create a table to store sales information
		CREATE TABLE IF NOT EXISTS sales (
			id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
			customer_id UUID REFERENCES customers(id) ON DELETE CASCADE,
-- 			sale_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
			 -- Snapshot of customer
			customer_name VARCHAR(255),
			customer_instagram_account VARCHAR(255),
			customer_phone BIGINT,
			customer_address VARCHAR(255),
			customer_city VARCHAR(255),
			customer_department VARCHAR(255),
			customer_comments VARCHAR(255),
			-- End Snapshot of customer
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		-- Create a mapping table to associate products with sales
		CREATE TABLE IF NOT EXISTS sale_products (
			sale_id UUID REFERENCES sales(id) ON DELETE CASCADE,
			product_variation_id UUID REFERENCES product_variations(id) ON DELETE CASCADE,
			PRIMARY KEY (sale_id, product_variation_id)
		);
    `)
	if err != nil {
		return err
	}

	var triggerExists bool
	// START product_variations trigger -----------------------------
	// Check if the updatedAt trigger already exists
	err = s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM pg_trigger 
         	WHERE tgname = 'product_variations_updated_at_trigger' 
           	AND tgrelid = 'product_variations'::regclass)
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
            
            CREATE TRIGGER product_variations_updated_at_trigger
            BEFORE UPDATE ON product_variations
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
            WHERE tgname = 'product_variations_created_at_trigger' 
            AND tgrelid = 'product_variations'::regclass)
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
            
            CREATE TRIGGER product_variations_created_at_trigger
            BEFORE INSERT ON product_variations
            FOR EACH ROW
            EXECUTE FUNCTION set_created_at();
        `)
		if err != nil {
			return err
		}
	}
	// END sales trigger -------------------------------

	// START sales trigger -----------------------------
	// Check if the updatedAt trigger already exists
	err = s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM pg_trigger 
         	WHERE tgname = 'sales_updated_at_trigger' 
           	AND tgrelid = 'sales'::regclass)
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
            
            CREATE TRIGGER sales_updated_at_trigger
            BEFORE UPDATE ON sales
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
            WHERE tgname = 'sales_created_at_trigger' 
            AND tgrelid = 'sales'::regclass)
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
            
            CREATE TRIGGER sales_created_at_trigger
            BEFORE INSERT ON sales
            FOR EACH ROW
            EXECUTE FUNCTION set_created_at();
        `)
		if err != nil {
			return err
		}
	}
	// END sales trigger -----------------------------

	return nil
}

func (s *PostgresStore) CreateSale(sale *SaleWithProducts) error {
	pvIDs, err := createProductVariations(sale, s)
	if err != nil {
		return err
	}

	customer, err := s.GetCustomerByID(sale.CustomerID)
	if err != nil {
		return err
	}

	saleID, err := createSale(customer, s)
	if err != nil {
		return err
	}

	err = createSaleProducts(saleID, pvIDs, s)
	if err != nil {
		return err
	}

	// Set the ID of the inserted sale (sales table)
	sale.ID = saleID

	return nil
}

// createProductVariations inserts product variations
func createProductVariations(sale *SaleWithProducts, s *PostgresStore) ([]string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("error rolling back transaction: %v", rollbackErr)
			}
			return
		}
		if err := tx.Commit(); err != nil {
			log.Printf("error committing transaction: %v", err)
		}
	}()

	copyIn, err := tx.Prepare(pq.CopyIn(
		"product_variations",
		"product_id",
		"color",
		"price",
		"created_at",
		"updated_at",
	))
	if err != nil {
		return nil, fmt.Errorf("error preparing COPY statement: %v", err)
	}
	defer func() {
		if err := copyIn.Close(); err != nil {
			log.Printf("error closing prepared statement: %v", err)
		}
	}()

	for _, product := range sale.Products {
		_, err := copyIn.Exec(
			product.ProductID,
			product.Color,
			product.Price,
			product.CreatedAt,
			product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error queuing data for COPY: %v", err)
		}
	}

	_, err = copyIn.Exec()
	if err != nil {
		return nil, fmt.Errorf("error executing COPY statement: %v", err)
	}

	// Fetch IDs of the inserted records
	rows, err := tx.Query("SELECT id FROM product_variations ORDER BY created_at DESC LIMIT $1", len(sale.Products))
	if err != nil {
		return nil, fmt.Errorf("error fetching IDs: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}()

	var insertedIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning ID: %v", err)
		}
		insertedIDs = append(insertedIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return insertedIDs, nil
}

func createSale(customerInSale *Customer, s *PostgresStore) (string, error) {
	query := `
        INSERT INTO sales (
			customer_id,
			customer_name,
			customer_instagram_account,
            customer_phone,
			customer_address,
			customer_city,
			customer_department,
			customer_comments,
		    created_at,
		    updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `

	// Resets customer's createdAt and updatedAt
	customerInSale.CreatedAt = time.Time{}
	customerInSale.UpdatedAt = time.Time{}

	var id string
	err := s.db.QueryRow(
		query,
		customerInSale.ID,
		customerInSale.Name,
		customerInSale.InstagramAccount,
		customerInSale.Phone,
		customerInSale.Address,
		customerInSale.City,
		customerInSale.Department,
		customerInSale.Comments,
		customerInSale.CreatedAt,
		customerInSale.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	// Set the ID of the inserted sale
	customerInSale.ID = id // sale ID

	return id, nil
}

func createSaleProducts(saleID string, pvIDs []string, s *PostgresStore) error {
	// Begin a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	// Prepare the COPY command
	copyIn, err := tx.Prepare(pq.CopyIn(
		"sale_products",
		"sale_id",
		"product_variation_id",
	))
	if err != nil {
		return fmt.Errorf("error preparing COPY statement: %v", err)
	}
	defer func(copyIn *sql.Stmt) {
		err := copyIn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(copyIn)

	// Queue the data to be copied
	for _, id := range pvIDs {
		_, err = copyIn.Exec(
			saleID,
			id,
		)
		if err != nil {
			return fmt.Errorf("error queuing data for COPY: %v", err)
		}
	}

	// Execute the COPY command
	_, err = copyIn.Exec()
	if err != nil {
		return fmt.Errorf("error executing COPY statement: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func (s *PostgresStore) GetSales() ([]*SaleResponse, error) {
	rows, err := s.db.Query(`
		SELECT
			s.id,
			s.customer_id,
			s.customer_name,
			s.customer_instagram_account,
			s.customer_phone,
			s.customer_address,
			s.customer_city,
			s.customer_department,
			s.customer_comments,
			s.created_at,
			s.updated_at,
			JSON_AGG(JSON_BUILD_OBJECT(
				'id', pv.id,
				'color', pv.color,
				'price', pv.price
			)) AS product_variations
		FROM
			sales s
		JOIN
			sale_products sp ON s.id = sp.sale_id
		JOIN
			product_variations pv ON sp.product_variation_id = pv.id
		GROUP BY
			s.id,
			s.customer_id,
			s.customer_name,
			s.customer_instagram_account,
			s.customer_phone,
			s.customer_address,
			s.customer_city,
			s.customer_department,
			s.customer_comments,
			s.created_at,
			s.updated_at;
`)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var sales []*SaleResponse
	for rows.Next() {
		sale, err := scanIntoSales(rows)
		if err != nil {
			return nil, err
		}

		sales = append(sales, sale)
	}

	return sales, nil
}

func scanIntoSales(rows *sql.Rows) (*SaleResponse, error) {
	sale := new(SaleResponse)
	var productVariationsJSON []byte
	err := rows.Scan(
		&sale.ID,
		&sale.CustomerID,
		&sale.CustomerName,
		&sale.CustomerInstagramAccount,
		&sale.CustomerPhone,
		&sale.CustomerAddress,
		&sale.CustomerCity,
		&sale.CustomerDepartment,
		&sale.CustomerComments,
		&sale.CreatedAt,
		&sale.UpdatedAt,
		&productVariationsJSON, // Scan JSON data into a []byte
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data into slice
	err = json.Unmarshal(productVariationsJSON, &sale.ProductVariations)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling product_variations JSON: %v", err)
	}

	return sale, nil
}

func (s *PostgresStore) GetSaleByID(id string) (*SaleResponse, error) {
	rows, err := s.db.Query(`
		SELECT
			s.id,
			s.customer_id,
			s.customer_name,
			s.customer_instagram_account,
			s.customer_phone,
			s.customer_address,
			s.customer_city,
			s.customer_department,
			s.customer_comments,
			s.created_at,
			s.updated_at,
			JSON_AGG(JSON_BUILD_OBJECT(
				'id', pv.id,
				'color', pv.color,
				'price', pv.price
			)) AS product_variations
		FROM
			sales s
		JOIN
			sale_products sp ON s.id = sp.sale_id
		JOIN
			product_variations pv ON sp.product_variation_id = pv.id
		WHERE
			s.id = $1
		GROUP BY
			s.id,
			s.customer_id,
			s.customer_name,
			s.customer_instagram_account,
			s.customer_phone,
			s.customer_address,
			s.customer_city,
			s.customer_department,
			s.customer_comments,
			s.created_at,
			s.updated_at;
	`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoSales(rows)
	}

	return nil, fmt.Errorf("sale [%s] not found", id)
}
