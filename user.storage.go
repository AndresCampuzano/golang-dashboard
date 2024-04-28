package main

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *PostgresStore) CreateUsersTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
            first_name VARCHAR(255) NOT NULL,
            last_name VARCHAR(255) NOT NULL,
            email VARCHAR(255) NOT NULL UNIQUE,
            encrypted_password VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateUser(user *User) error {
	query := `
        INSERT INTO users (first_name, last_name, email, encrypted_password, created_at) 
        VALUES ($1, $2, $3, $4, $5) 
        RETURNING id
    `

	var id string
	err := s.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.EncryptedPassword,
		user.CreatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}

	// Set the ID of the inserted user
	user.ID = id

	return nil
}

func (s *PostgresStore) GetUsers() ([]*User, error) {
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var users []*User
	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *PostgresStore) GetUserByID(id string) (*User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user [%s] not found", id)
}

func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user [%s] not found", email)
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.EncryptedPassword,
		&user.CreatedAt,
	)

	return user, err
}
