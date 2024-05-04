package main

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *PostgresStore) CreateExpensesTable() error {
	// Create the table if it doesn't exist
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS expenses (
            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            price BIGINT NOT NULL,
            type VARCHAR(255) NOT NULL,
            description VARCHAR(255) NOT NULL,
            currency VARCHAR(255) NOT NULL,
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
         	WHERE tgname = 'expenses_updated_at_trigger' 
           	AND tgrelid = 'expenses'::regclass)
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
            
            CREATE TRIGGER expenses_updated_at_trigger
            BEFORE UPDATE ON expenses
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
            WHERE tgname = 'expenses_created_at_trigger' 
            AND tgrelid = 'expenses'::regclass)
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
            
            CREATE TRIGGER expenses_created_at_trigger
            BEFORE INSERT ON expenses
            FOR EACH ROW
            EXECUTE FUNCTION set_created_at();
        `)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStore) CreateExpense(expense *Expense) error {
	query := `
        INSERT INTO expenses (
			name, 
			price, 
			type, 
			description, 
			currency, 
			created_at, 
			updated_at 
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) 
        RETURNING id
    `

	var id string
	err := s.db.QueryRow(
		query,
		expense.Name,
		expense.Price,
		expense.Type,
		expense.Description,
		expense.Currency,
		expense.CreatedAt,
		expense.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}

	// Set the ID of the inserted expense
	expense.ID = id

	return nil
}

func (s *PostgresStore) GetExpenseByID(id string) (*Expense, error) {
	rows, err := s.db.Query("SELECT * FROM expenses WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoExpenses(rows)
	}

	return nil, fmt.Errorf("expense [%s] not found", id)
}

func scanIntoExpenses(rows *sql.Rows) (*Expense, error) {
	expense := new(Expense)
	err := rows.Scan(
		&expense.ID,
		&expense.Name,
		&expense.Price,
		&expense.Type,
		&expense.Description,
		&expense.Currency,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	return expense, err
}

func (s *PostgresStore) GetExpenses() ([]*Expense, error) {
	rows, err := s.db.Query("SELECT * FROM expenses")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var expenses []*Expense
	for rows.Next() {
		expense, err := scanIntoExpenses(rows)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func (s *PostgresStore) UpdateExpense(expense *Expense) error {
	query := `
		UPDATE expenses
		SET 
		    name = $1, 
		    price = $2, 
		    type = $3, 
		    description = $4,
		    currency = $5
		WHERE id = $6
	`

	_, err := s.db.Exec(
		query,
		expense.Name,
		expense.Price,
		expense.Type,
		expense.Description,
		expense.Currency,
		expense.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteExpense(id string) error {
	_, err := s.db.Exec("DELETE FROM expenses WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
