package store

import "book-tracker/models"

type BookStore interface {
	CreateBook() error
	GetBook() (*models.Book, error)
	ListBooks() ([]*models.Book, error)
	UpdateBook() error
	DeleteBook() error
	GetStats() (totalRead, readingProcess int, popularAuthor string, err error)
}
