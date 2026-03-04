package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type Database struct {
	db *sql.DB
}

func (d *Database) Close() error {
	return d.db.Close()
}

func GetDatabase(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) CreateBookTable() error {
	insertQuery := `CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT,
		openID TEXT UNIQUE NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err := d.db.Exec(insertQuery)
	if err != nil {
		return fmt.Errorf("Error creating table: %w", err)
	}

	return nil
}

func (d *Database) CreateUserTable() error {
	createQuery := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);`

	_, err := d.db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("Error creating user table: %w", err)
	}

	return nil
}

func (d *Database) AddBook(title, author, openID string, user_id int) (Book, error) {
	if strings.TrimSpace(title) == "" {
		return Book{}, errors.New("Book title cannot be empty")
	}
	if strings.TrimSpace(author) == "" {
		return Book{}, errors.New("Book author cannot be empty.")
	}
	query := `INSERT INTO books (title, author, openID, user_id) VALUES (?, ?, ?, ?);`

	result, insertErr := d.db.Exec(query, title, author, openID, user_id)

	if insertErr != nil {
		return Book{}, insertErr
	}

	id, insertErr := result.LastInsertId()
	if insertErr != nil {
		return Book{}, insertErr
	}
	return Book{
		ID:      id,
		Title:   title,
		Author:  author,
		OpenID:  openID,
		User_id: user_id,
	}, nil
}

func (d *Database) AddUser(name string) (int64, error) {
	query := "INSERT INTO users (name) values (?);"

	result, insertErr := d.db.Exec(query, name)

	if insertErr != nil {
		return 0, insertErr
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
