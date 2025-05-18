package store

import "book-tracker/models"

type BookStore interface {
	CreateBook(book *models.Book) error
	GetBook(id string) (*models.Book, error)
	ListBooks(status string, limit, offset int, title, author string) ([]*models.Book, error)
	UpdateBook(book *models.Book) error
	DeleteBook(id string) error
	GetStats() (totalRead, readingProcess int, popularAuthor string, err error)
	CountBooks() (total int, byStatus map[string]int)
}
