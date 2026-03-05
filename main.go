package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func (d *Database) SearchHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")

	results, err := d.SearchBooks(term)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, results)
}

func (d *Database) AddBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = d.AddUser(newUser.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database errro: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s inserted successfully.", newUser.Username)
}

func main() {
	db, err := sql.Open("sqlite3", "/data/book_database.db")

	app := &Database{db: db}
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Database connection successful!")

	http.HandleFunc("/", handler)
	http.HandleFunc("/search", app.SearchHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
