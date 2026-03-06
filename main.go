package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	DB   Database
	Tmpl *template.Template
}

func (a *App) IndexHandler(w http.ResponseWriter, r *http.Request) {
	a.Tmpl.Execute(w, nil)
}

func (a *App) SearchHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")

	results, err := a.DB.SearchBooks(term)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, results)
}

func (a *App) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the form
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	err := a.DB.DeleteBook(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) AddBookHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = a.DB.AddUser(newUser.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("App errro: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s inserted successfully.", newUser.Username)
}

func main() {
	dbConn, err := sql.Open("sqlite3", "./book_database.db")

	if err != nil {
		log.Fatal(err)
	}
	db := Database{db: dbConn}
	defer db.Close()

	db.CreateBookTable()
	db.CreateUserTable()
	fmt.Println("Database connection successful!")

	tmpl := template.Must(template.ParseFiles("index.html"))

	app := App{DB: db, Tmpl: tmpl}

	http.HandleFunc("/", app.IndexHandler)
	http.HandleFunc("/search", app.SearchHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
