package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	db, err := sql.Open("sqlite3", "/data/book_database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Database connection successful!")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
