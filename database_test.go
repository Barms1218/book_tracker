package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
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
		name    string
		title   string
		author  string
		openID  string
		user_id int
		wantErr bool
	}

	addTestCases := []testCase{
		{
			name:    "Standard book insert",
			title:   "The Go Programming Language",
			author:  "Alan Donovan",
			openID:  "12345",
			user_id: 1,
			wantErr: false,
		},
		{
			name:    "Duplicate OpenID",
			title:   "Red Rising",
			author:  "Pierce Brown",
			openID:  "12345",
			user_id: 1,
			wantErr: true,
		},
		{
			name:    "Empty title",
			title:   "",
			author:  "Nobody",
			openID:  "54321",
			user_id: 1,
			wantErr: true,
		},
		{
			name:    "No User Found",
			title:   "The Hobbit",
			author:  "J.R.R Tolkien",
			openID:  "12345",
			user_id: 1,
			wantErr: true,
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

			if tc.name == "Duplicate OpenID" {
				_, err := db.AddBook(tc.title, tc.author, tc.openID, tc.user_id)

				if err != nil {
					t.Errorf("Setup insertion failed, enountered error : %v", err)
				}
			}

			book, err := db.AddBook(tc.title, tc.author, tc.openID, tc.user_id)

			if tc.wantErr != (err != nil) {
				t.Fatalf("Expected error: %v, got %v", tc.wantErr, err)
			}

			if !tc.wantErr {
				if book.ID != 1 {
					t.Errorf("Expected book ID of 1, got %d", book.ID)
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
		expectedTitles []string
	}
	getTestCases := []testCase{
		{
			name:       "All Results Contain King",
			searchTerm: "King",
			expectedTitles: []string{
				"The Return of the King",
				"Kingfisher",
				"The Shining",
				"The Wicked King",
			},
		},
		{
			name:           "Search With Gibberish",
			searchTerm:     "Xyjasdfkl",
			expectedTitles: []string{},
		},
		{
			name:       "Limit Results to 10",
			searchTerm: "The",
			expectedTitles: []string{
				"The Hobbit",
				"The Fellowship of the Ring",
				"The Two Towers",
				"The Return of the King",
				"The Shining",
				"The Wicked King",
				"Theo of Golden",
				"The Intruder",
				"The Giver",
				"The Friend of the Family",
				"The Will of the Many",
				"The Light of All That Falls",
			},
		},
	}
}
