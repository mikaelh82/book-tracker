package models

type BookStatus string

const (
	BookUnread   BookStatus = "unread"
	BookReading  BookStatus = "reading"
	BookComplete BookStatus = "complete"
)

type Book struct {
	ID     string
	Title  string
	Author string
	Status BookStatus
}
