package main

type Book struct {
	ID            int64
	Title         string
	Author        string
	OpenID        string
	ReadingStatus string
	User_id       int64
}

type User struct {
	ID       int64
	Username string
}

type JournalEntries struct {
	ID      int64
	book_id int64
	content string
}
