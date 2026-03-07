package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type App struct {
	DB   Database
	Tmpl *template.Template
}

type SearchPage struct {
	SearchResults  []Book
	SavedBooks     []Book
	CollectionSize int
}

func (a *App) IndexHandler(w http.ResponseWriter, r *http.Request) {
	localResults, _ := a.DB.SearchUserBooks(1)

	log.Printf("DEBUG: Loaded %d books for the 'My Collection' section", len(localResults))

	data := SearchPage{
		SearchResults:  nil,
		SavedBooks:     localResults,
		CollectionSize: len(localResults),
	}

	a.Tmpl.Execute(w, data)
}

func (a *App) SearchHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	if term == "" {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}

	results, err := FetchFromOpenLibrary(term)

	if err != nil {
		log.Printf("Open Library error %v", err)
		http.Error(w, "Failed to fetch books from Open Library", http.StatusInternalServerError)
	}

	localResults, _ := a.DB.SearchUserBooks(1)
	data := SearchPage{
		SearchResults:  results,
		SavedBooks:     localResults,
		CollectionSize: len(localResults),
	}

	err = a.Tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template error: %v", err)
	}
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	title := r.PostFormValue("title")
	author := r.PostFormValue("author")
	openID := r.PostFormValue("openid")

	log.Printf("Handler received: Title=%s, Author=%s, OpenID=%s", title, author, openID)
	_, err := a.DB.AddBook(title, author, openID, 1)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		log.Printf("Database insertion error: %v", err)
		http.Error(w, "Failed to save book due to error: %v.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func FetchFromOpenLibrary(query string) ([]Book, error) {
	url := fmt.Sprintf("https://openlibrary.org/search.json?q=%s&fields=title,author_name,key&limit=10", url.QueryEscape(query))
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	// Create a struct that will match fields with the ones returned by the open library quest
	var result struct {
		Docs []struct {
			Title  string   `json:"title"`
			Author []string `json:"author_name"`
			Key    string   `json:"key"`
		} `json:"docs"`
	}
	// Rate limits are 100 requests per IP every 5 minutes

	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	var books []Book
	for _, doc := range result.Docs {
		author := "Unknown"
		if len(doc.Author) > 0 {
			author = doc.Author[0]
		}
		books = append(books, Book{
			Title:  doc.Title,
			Author: author,
			OpenID: doc.Key,
		})
	}
	return books, nil

}

func main() {
	dbConn, err := sql.Open("sqlite3", "./book_database.db")

	if err != nil {
		log.Fatal(err)
	}
	db := Database{db: dbConn}
	defer db.Close()

	_, _ = dbConn.Exec("PRAGMA foreign_keys = ON;")

	db.CreateBookTable()
	db.CreateUserTable()

	db.AddUser("Branden")

	fmt.Println("Database connection successful!")

	tmpl := template.Must(template.ParseFiles("index.html"))

	app := App{DB: db, Tmpl: tmpl}

	http.HandleFunc("/", app.IndexHandler)
	http.HandleFunc("/search", app.SearchHandler)
	http.HandleFunc("/delete", app.DeleteBookHandler)
	http.HandleFunc("/add", app.AddBookHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
