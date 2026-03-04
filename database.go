package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func GetDatabase(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) CreateBookTable() (int, error) {
	insertQuery := `CREATE TABLE book (
		id INTEGER PRIMARY KEY,
		openID TEXT UNIQUE NOT NUll,
		title TEXT UNIQUE NOT NULL,
		author TEXT,
	);`

	id, err := d.db.Exec(insertQuery)
	if err != nil {
		return fmt.Errorf("Error creating table: %w", err)
	}

	return id, nil
}

func (d *Database) AddBook(title, author, openID string) (Book, error) {
	query := `INSERT INTO books (title, author, openID) VALUES (?, ?, ?);`

	result, insertErr := d.db.Exec(query, title, author, openID)

	if insertErr != nil {
		return Book{}, insertErr
	}

	id, insertErr := result.LastInsertId()
	if insertErr != nil {
		return Book{}, insertErr
	}
	return Book{
		id:     id,
		title:  title,
		author: author,
		openID: openID,
	}, nil
}
