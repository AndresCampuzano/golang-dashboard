-- Recreate primary keys
ALTER TABLE product_variations ADD PRIMARY KEY (id);
ALTER TABLE sales ADD PRIMARY KEY (id);
ALTER TABLE sale_products ADD PRIMARY KEY (sale_id, product_variation_id);
ALTER TABLE products ADD PRIMARY KEY (id);
ALTER TABLE customers ADD PRIMARY KEY (id);
ALTER TABLE expenses ADD PRIMARY KEY (id);
ALTER TABLE users ADD PRIMARY KEY (id);

-- Recreate foreign key constraints
ALTER TABLE product_variations ADD CONSTRAINT fk_product_variation_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
ALTER TABLE sales ADD CONSTRAINT fk_sale_customer_id FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE;
ALTER TABLE sale_products ADD CONSTRAINT fk_sale_products_sale_id FOREIGN KEY (sale_id) REFERENCES sales(id) ON DELETE CASCADE;
ALTER TABLE sale_products ADD CONSTRAINT fk_sale_products_product_variation_id FOREIGN KEY (product_variation_id) REFERENCES product_variations(id) ON DELETE CASCADE;


-- Recreate additional defaults
ALTER TABLE product_variations ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE product_variations ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE sales ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE sales ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE products ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE products ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE customers ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE customers ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE expenses ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE expenses ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE users ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;

-- Set NOT NULL constraints for products table
ALTER TABLE products ALTER COLUMN name SET NOT NULL;
ALTER TABLE products ALTER COLUMN price SET NOT NULL;
ALTER TABLE products ALTER COLUMN image SET NOT NULL;

-- Set NOT NULL constraints for customers table
ALTER TABLE customers ALTER COLUMN name SET NOT NULL;
ALTER TABLE customers ALTER COLUMN instagram_account SET NOT NULL;
ALTER TABLE customers ALTER COLUMN address SET NOT NULL;
ALTER TABLE customers ALTER COLUMN city SET NOT NULL;
ALTER TABLE customers ALTER COLUMN department SET NOT NULL;
ALTER TABLE customers ALTER COLUMN comments SET NOT NULL;
ALTER TABLE customers ALTER COLUMN cc SET NOT NULL;

-- Set NOT NULL constraints for expenses table
ALTER TABLE expenses ALTER COLUMN name SET NOT NULL;
ALTER TABLE expenses ALTER COLUMN price SET NOT NULL;
ALTER TABLE expenses ALTER COLUMN type SET NOT NULL;
ALTER TABLE expenses ALTER COLUMN description SET NOT NULL;
ALTER TABLE expenses ALTER COLUMN currency SET NOT NULL;

-- Set NOT NULL constraints for users table
ALTER TABLE users ALTER COLUMN first_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN last_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ALTER COLUMN encrypted_password SET NOT NULL;

-- Set NOT NULL constraints for product_variations table
ALTER TABLE product_variations ALTER COLUMN color SET NOT NULL;
ALTER TABLE product_variations ALTER COLUMN price SET NOT NULL;

-- Set NOT NULL constraints for sales table
ALTER TABLE sales ALTER COLUMN customer_name SET NOT NULL;
ALTER TABLE sales ALTER COLUMN customer_address SET NOT NULL;
ALTER TABLE sales ALTER COLUMN customer_address SET NOT NULL;
ALTER TABLE sales ALTER COLUMN customer_city SET NOT NULL;
ALTER TABLE sales ALTER COLUMN customer_department SET NOT NULL;
ALTER TABLE sales ALTER COLUMN customer_comments SET NOT NULL;
ALTER TABLE sales ALTER COLUMN customer_cc SET NOT NULL;

-- Set NOT NULL constraint for users table
ALTER TABLE users ALTER COLUMN email SET NOT NULL;

-- Add UNIQUE constraint for email column in users table
ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);

-- Set DEFAULT uuid_generate_v4() for id columns in all tables
ALTER TABLE product_variations ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE sales ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE products ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE customers ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE expenses ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4();

