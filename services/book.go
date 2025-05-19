package services

import (
	"book-tracker/models"
	"book-tracker/store"
	"context"
)

type BookService interface {
	CreateBook(ctx context.Context, book *models.Book) error
	GetBook(ctx context.Context, id string) (*models.Book, error)
	ListBooks(ctx context.Context, status string, limit, offset int) ([]*models.Book, error)
	UpdateBook(ctx context.Context, book *models.Book) error
	DeleteBook(ctx context.Context, id string) error
}

type bookService struct {
	store store.BookStore
}

func NewBookService(store store.BookStore) BookService {
	return &bookService{store: store}
}

func (s *bookService) CreateBook(ctx context.Context, book *models.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	if err := book.GenerateID(); err != nil {
		return err
	}
	return s.store.CreateBook(ctx, book)
}

func (s *bookService) GetBook(ctx context.Context, id string) (*models.Book, error) {
	return s.store.GetBook(ctx, id)
}

func (s *bookService) ListBooks(ctx context.Context, status string, limit, offset int) ([]*models.Book, error) {
	// NOTE:
	// We can extend functionality later on author and book filtering if we wish
	// but lets keep these neutral for now and only adapt for pagination or infinite-scroll (lets keep status)
	return s.store.ListBooks(ctx, status, limit, offset, "", "")
}

func (s *bookService) UpdateBook(ctx context.Context, book *models.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	return s.store.UpdateBook(ctx, book)
}

func (s *bookService) DeleteBook(ctx context.Context, id string) error {
	return s.store.DeleteBook(ctx, id)
}
