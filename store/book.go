package store

import (
	"book-tracker/models"
	"context"
)

type BookStore interface {
	CreateBook(ctx context.Context, book *models.Book) error
	GetBook(ctx context.Context, id string) (*models.Book, error)
	ListBooks(ctx context.Context, status string, limit, offset int, title, author string) ([]*models.Book, error)
	UpdateBook(ctx context.Context, book *models.Book) error
	DeleteBook(ctx context.Context, id string) error
	GetStats(ctx context.Context) (totalRead, readingProcess int, popularAuthor string, err error)
	CountBooks(ctx context.Context) (total int, byStatus map[string]int)
}
