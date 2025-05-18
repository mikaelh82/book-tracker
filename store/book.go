package store

import (
	"book-tracker/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type BookStore interface {
	CreateBook(ctx context.Context, book *models.Book) error
	GetBook(ctx context.Context, id string) (models.Book, error)
	ListBooks(ctx context.Context, status string, limit, offset int, title, author string) ([]*models.Book, error)
	UpdateBook(ctx context.Context, book *models.Book) error
	DeleteBook(ctx context.Context, id string) error
	GetStats(ctx context.Context) (totalRead, readingProcess int, popularAuthor string, err error)
	CountBooks(ctx context.Context) (total int, byStatus map[string]int)
}

type bookStore struct {
	db *sql.DB
}

func NewBookStore(db *sql.DB) BookStore {
	return &bookStore{db: db}
}

func (s *bookStore) CreateBook(ctx context.Context, book *models.Book) error {
	_, err := s.db.ExecContext(ctx, `
	INSERT INTO books (id, title, author, status) VALUES (?, ?, ?, ?)`,
		book.ID, book.Title, book.Author, book.Status)

	// TODO: Could be more explicit here and add an error type if the book already exsists given that (title, author) is unique

	if err != nil {
		return fmt.Errorf("insert book: %w", err)
	}

	return nil
}

func (s *bookStore) GetBook(ctx context.Context, id string) (models.Book, error) {

	var book models.Book

	err := s.db.QueryRowContext(ctx, `
	SELECT id, title, author, status
	FROM books
	WHERE id = ?`, id).Scan(&book.ID, &book.Title, &book.Author, &book.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Book{}, ErrBookNotFound
		}
		return models.Book{}, fmt.Errorf("get book: %w", err)
	}

	return book, nil
}

func (s *bookStore) ListBooks(ctx context.Context, status string, limit, offset int, title, author string) ([]*models.Book, error) {
	// TODO: unimplemented
	return nil, nil
}

func (s *bookStore) UpdateBook(ctx context.Context, book *models.Book) error {
	// TODO: unimplemented
	return nil
}

func (s *bookStore) DeleteBook(ctx context.Context, id string) error {
	// TODO: unimplemented
	return nil
}

func (s *bookStore) GetStats(ctx context.Context) (totalRead, readingProgress int, popularAuthor string, err error) {
	// TODO: unimplemented
	return 0, 0, "", nil
}

func (s *bookStore) CountBooks(ctx context.Context) (total int, byStatus map[string]int) {
	// TODO: unimplemented
	return 0, nil
}
