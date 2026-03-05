package main

type Book struct {
	ID      int64
	Title   string
	Author  string
	OpenID  string
	User_id int64
}

type User struct {
	ID       int64
	Username string
}
