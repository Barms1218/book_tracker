package main

type Book struct {
	ID      int64
	Title   string
	Author  string
	OpenID  string
	User_id int
}

type User struct {
	ID       int64
	Username string
}
