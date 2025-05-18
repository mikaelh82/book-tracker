package store

import (
	"context"
	"database/sql"
	"fmt"
)

type StatsStore interface {
	GetStats(ctx context.Context) (totalRead, readingProgress int, popularAuthor string, err error)
}

type statsStore struct {
	db *sql.DB
}

func NewStatsStore(db *sql.DB) StatsStore {
	return &statsStore{db: db}
}

func (s *statsStore) GetStats(ctx context.Context) (totalRead, readingProgress int, popularAuthor string, err error) {
	err = s.db.QueryRowContext(ctx, `
	SELECT COUNT(*)
	FROM books
	WHERE status = 'complete'
	`).Scan(&totalRead)

	if err != nil {
		return 0, 0, "", fmt.Errorf("query total completed: %w", err)
	}

	var totalBooks float64

	err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM books").Scan(&totalBooks)
	if err != nil {
		return 0, 0, "", fmt.Errorf("query total books: %w", err)
	}

	var readingBooks float64
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM books
		WHERE status = 'reading'
	`).Scan(&readingBooks)
	if err != nil {
		return 0, 0, "", fmt.Errorf("query reading books: %w", err)
	}

	if totalBooks > 0 {
		readingProgress = int((readingBooks / totalBooks) * 100)
	}

	err = s.db.QueryRowContext(ctx, `
		SELECT author 
		FROM books 
		GROUP BY author 
		ORDER BY COUNT(*) DESC 
		LIMIT 1
	`).Scan(&popularAuthor)
	if err == sql.ErrNoRows {
		popularAuthor = "N/A"
	} else if err != nil {
		return 0, 0, "", fmt.Errorf("query popular author: %w", err)
	}

	return totalRead, readingProgress, popularAuthor, nil

}
