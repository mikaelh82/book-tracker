package store

import (
	"book-tracker/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type BookStore interface {
	CreateBook(ctx context.Context, book *models.Book) error
	GetBook(ctx context.Context, id string) (*models.Book, error)
	ListBooks(ctx context.Context, status string, limit, offset int, title, author string) ([]*models.Book, error)
	UpdateBook(ctx context.Context, book *models.Book) error
	DeleteBook(ctx context.Context, id string) error
	CountBooks(ctx context.Context) (total int, byStatus map[string]int, err error)
}

type bookStore struct {
	db *sql.DB
}

func NewBookStore(db *sql.DB) BookStore {
	return &bookStore{db: db}
}

func (s *bookStore) CreateBook(ctx context.Context, book *models.Book) error {
	// NOTE: Documentation: https://pkg.go.dev/database/sql#Conn.ExecContext
	_, err := s.db.ExecContext(ctx, `
        INSERT INTO books (id, title, author, status) VALUES (?, ?, ?, ?)`,
		book.ID, book.Title, book.Author, book.Status)
	if err != nil {
		return fmt.Errorf("create book: %w", err)
	}
	return nil
}

func (s *bookStore) GetBook(ctx context.Context, id string) (*models.Book, error) {
	// NOTE: Documentation: https://pkg.go.dev/database/sql#DB.QueryRowContext
	var book models.Book
	err := s.db.QueryRowContext(ctx, `
        SELECT id, title, author, status
        FROM books
        WHERE id = ?`, id).Scan(&book.ID, &book.Title, &book.Author, &book.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("get book: %w", err)
	}
	return &book, nil // Lets stress test Garbage Collector =)
}

// IMPORTANT:
// This can initially be seen as very unsafe and thats is indeed true but keep in mind that we handle
// all validations in the service layer so nothing dangerous will be injected into here
// Also, in this case, i am careful not to bring in too many external libraries but this could be simplified
// alot with a ORM like Prisma (or GORM of go in this case). But that also adds overhead
func (s *bookStore) ListBooks(ctx context.Context, status string, limit, offset int, title, author string) ([]*models.Book, error) {
	// NOTE: Documentation: https://pkg.go.dev/database/sql#DB.QueryContext
	query := "SELECT id, title, author, status FROM books"
	args := []any{}
	conditions := []string{}
	if status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}
	if title != "" {
		conditions = append(conditions, "title = ?")
		args = append(args, title)
	}
	if author != "" {
		conditions = append(conditions, "author = ?")
		args = append(args, author)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY title ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query books: %w", err)
	}
	defer rows.Close()

	books := []*models.Book{}
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Status); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, &book) // Lets stress test Garbage Collector =)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return books, nil
}

func (s *bookStore) UpdateBook(ctx context.Context, book *models.Book) error {
	// NOTE: Documentation: https://pkg.go.dev/database/sql#DB.ExecContext
	result, err := s.db.ExecContext(ctx, `
        UPDATE books
        SET title = ?, author = ?, status = ?
        WHERE id = ?
    `, book.Title, book.Author, book.Status, book.ID)
	if err != nil {
		return fmt.Errorf("update book: %w", err)
	}
	// NOTE: Documentation: https://pkg.go.dev/database/sql#Result
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrBookNotFound
	}
	return nil
}

func (s *bookStore) DeleteBook(ctx context.Context, id string) error {
	// NOTE: Documentation: https://pkg.go.dev/database/sql#DB.ExecContext
	result, err := s.db.ExecContext(ctx, "DELETE FROM books WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrBookNotFound
	}
	return nil
}

func (s *bookStore) CountBooks(ctx context.Context) (total int, byStatus map[string]int, err error) {
	// NOTE: Documentation: https://pkg.go.dev/database/sql#DB.QueryContext
	// NOTE:
	// var total int <--- Not possible (or already in-scope) as it somehow declares the variables in the function
	// return declaration as in-scope variables?
	err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM books").Scan(&total)
	if err != nil {
		return 0, nil, fmt.Errorf("count total books: %w", err)
	}
	byStatus = make(map[string]int)
	rows, err := s.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM books GROUP BY status")
	if err != nil {
		return total, nil, fmt.Errorf("query books by status: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return total, byStatus, fmt.Errorf("scan status count: %w", err)
		}
		byStatus[status] = count
	}
	if err := rows.Err(); err != nil {
		return total, byStatus, fmt.Errorf("rows error: %w", err)
	}
	return total, byStatus, nil
}
