package models

import (
	"strings"
)

type BookStatus string

const (
	BookUnread   BookStatus = "unread"
	BookReading  BookStatus = "reading"
	BookComplete BookStatus = "complete"
)

type Book struct {
	ID     string     `json:"id"`
	Title  string     `json:"title"`
	Author string     `json:"author"`
	Status BookStatus `json:"status"`
}

func (b *Book) Validate() error {
	b.Title = strings.TrimSpace(b.Title)
	b.Author = strings.TrimSpace(b.Author)

	trimmedStatus := strings.TrimSpace(string(b.Status))
	lowercaseStatus := strings.ToLower(trimmedStatus)
	status := BookStatus(lowercaseStatus)

	if b.Title == "" {
		//ERROR MISSING TITLE
	}

	if b.Author == "" {
		// ERROR MISSING AUTHOR
	}

	switch status {
	case BookUnread, BookReading, BookComplete:
		// ALL GOOD
	default:
		// ERROR UNKNOWN STATUS
	}

	return nil
}
