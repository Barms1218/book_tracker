package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "/data/book_database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Database connection successful!")
}
