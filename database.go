package main

import (
	"database/sql"
	"fmt"
	"os"
)

func NewPostgresStore() (*PostgresStore, error) {
	// Retrieve environment variables
	user := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB_NAME")
	password := os.Getenv("POSTGRES_PASSWORD")

	// Construct connection string
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", user, dbName, password)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connectivity
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Install UUID on postgres
	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		fmt.Println("Error creating uuid-ossp extension:", err)
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}
