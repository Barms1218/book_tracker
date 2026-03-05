package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"testing"
)

func SetupTestDB(t *testing.T) *Database {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory db: %v", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		t.Fatalf("Failed to enable foreign keys: %v", err)
	}

	d := &Database{db: db}

	if err := d.CreateUserTable(); err != nil {
		t.Fatalf("Failed to create user table: %v", err)
	}

	if err := d.CreateBookTable(); err != nil {
		t.Fatalf("Failed to create book table: %v", err)
	}

	return d
}

func TestAddBook(t *testing.T) {
	type testCase struct {
		name        string
		title       string
		author      string
		openuser_id string
		user_id     int64
		wantErr     bool
	}

	addTestCases := []testCase{
		{
			name:        "Standard book insert",
			title:       "The Go Programming Language",
			author:      "Alan Donovan",
			openuser_id: "12345",
			user_id:     1,
			wantErr:     false,
		},
		{
			name:        "Duplicate Openuser_id",
			title:       "Red Rising",
			author:      "Pierce Brown",
			openuser_id: "12345",
			user_id:     1,
			wantErr:     true,
		},
		{
			name:        "Empty title",
			title:       "",
			author:      "Nobody",
			openuser_id: "54321",
			user_id:     1,
			wantErr:     true,
		},
		{
			name:        "No User Found",
			title:       "The Hobbit",
			author:      "J.R.R Tolkien",
			openuser_id: "12345",
			user_id:     1,
			wantErr:     true,
		},
	}

	for _, tc := range addTestCases {
		t.Run(tc.name, func(t *testing.T) {
			db := SetupTestDB(t)
			defer db.Close()

			// Create test user
			if tc.name != "No User Found" {
				db.AddUser("Test User")
			}

			if tc.name == "Duplicate Openuser_id" {
				_, err := db.AddBook(tc.title, tc.author, tc.openuser_id, tc.user_id)

				if err != nil {
					t.Errorf("Setup insertion failed, enountered error : %v", err)
				}
			}

			book, err := db.AddBook(tc.title, tc.author, tc.openuser_id, tc.user_id)

			if tc.wantErr != (err != nil) {
				t.Fatalf("Expected error: %v, got %v", tc.wantErr, err)
			}

			if !tc.wantErr {
				if book.User_id != 1 {
					t.Errorf("Expected book user_id of 1, got %d", book.ID)
				}
				if book.Title != tc.title {
					t.Errorf("Expected title %s, got %v", tc.title, book.Title)
				}
			}
		})
	}
}

func GetBookByTitleTest(t *testing.T) {
	type testCase struct {
		name           string
		searchTerm     string
		seededBooks    []Book
		expectedLength int
	}
	getTestCases := []testCase{
		{
			name:       "All Results Contain King",
			searchTerm: "King",
			seededBooks: []Book{
				Book{Title: "The Return of the King", Author: "J.R.R Tolkien", OpenID: "12345"},
				Book{Title: "Kingfisher", Author: "Test Author", OpenID: "67890"},
				Book{Title: "The Shinig", Author: "Stephen King", OpenID: "54321"},
				Book{Title: "The Wicked King", Author: "Holly Black", OpenID: "09876"},
			},
			expectedLength: 4,
		},
		{
			name:       "Search With Gibberish",
			searchTerm: "Xyjasdfkl",
			seededBooks: []Book{
				Book{Title: "One", Author: "One", OpenID: "One"},
				Book{Title: "Two", Author: "Two", OpenID: "Two"},
				Book{Title: "Three", Author: "Three", OpenID: "Three"},
			},
			expectedLength: 0,
		},
		{
			name:       "Limit Results to 10",
			searchTerm: "The",
			seededBooks: []Book{
				Book{Title: "The First", Author: "A", OpenID: "1"},
				Book{Title: "The Second", Author: "A", OpenID: "2"},
				Book{Title: "The Third", Author: "A", OpenID: "3"},
				Book{Title: "The Fourth", Author: "A", OpenID: "4"},
				Book{Title: "The Fifth", Author: "A", OpenID: "5"},
				Book{Title: "The Sixth", Author: "A", OpenID: "6"},
				Book{Title: "The Seventh", Author: "A", OpenID: "7"},
				Book{Title: "The Eighth", Author: "A", OpenID: "8"},
				Book{Title: "The Ninth", Author: "A", OpenID: "9"},
				Book{Title: "The Tenth", Author: "A", OpenID: "10"},
				Book{Title: "The Eleventh", Author: "A", OpenID: "11"},
				Book{Title: "The Twelfth", Author: "A", OpenID: "12"},
			},
			expectedLength: 10,
		},
	}
	for _, tc := range getTestCases {
		t.Run(tc.name, func(t *testing.T) {
			db := SetupTestDB(t)
			defer db.Close()

			user_id, err := db.AddUser("Test User")
			if err != nil {
				t.Fatalf("Error encountered adding user: %v", err)
			}

			for _, book := range tc.seededBooks {
				db.AddBook(book.Title, book.Author, book.OpenID, user_id)
			}

			books, err := db.SearchBooks(tc.searchTerm)
			if err != nil {
				t.Fatalf("Error encountered searching for books with term %s: %v", tc.searchTerm, err)
			}

			if len(books) != tc.expectedLength {
				t.Fatalf("Search function was meant to return %d books, returned %d", tc.expectedLength, len(books))
			}

			for _, book := range books {
				hasTitle := strings.Contains(strings.ToLower(book.Title), strings.ToLower(tc.searchTerm))
				hasAuthor := strings.Contains(strings.ToLower(book.Author), strings.ToLower(tc.searchTerm))
				if !hasTitle && !hasAuthor {
					t.Fatalf("Book not containing search term was returned.")
				}
			}
		})
	}
}
