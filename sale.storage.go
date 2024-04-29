package main

func (s *PostgresStore) CreateSalesTablesWithRelations() error {
	// Create the table if it doesn't exist
	_, err := s.db.Exec(`
        -- Create a table to store product variations
		CREATE TABLE IF NOT EXISTS product_variations (
			id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
			product_id UUID REFERENCES products(id), -- Foreign key referencing the products table
			color VARCHAR(20) NOT NULL,
			price BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		-- Create a table to store sales information
		CREATE TABLE IF NOT EXISTS sales (
			id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
			customer_id UUID REFERENCES customers(id), -- Foreign key referencing the customers table
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
			sale_id UUID REFERENCES sales(id),
			product_variation_id UUID REFERENCES product_variations(id),
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
