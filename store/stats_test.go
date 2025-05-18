package store

import (
	"context"
	"testing"

	"book-tracker/models"

	"github.com/google/uuid"
)

func TestStatsStore(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	store := NewStatsStore(db)
	ctx := context.Background()

	addBook := func(id string, title string, author string, status models.BookStatus) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO books (id, title, author, status) VALUES (?, ?, ?, ?)`,
			id, title, author, status)
		if err != nil {
			t.Fatalf("Failed to insert book: %v", err)
		}
	}

	t.Run("EmptyDatabase", func(t *testing.T) {
		totalRead, readingProgress, popularAuthor, err := store.GetStats(ctx)
		if err != nil {
			t.Errorf("GetStats failed: %v", err)
		}
		if totalRead != 0 {
			t.Errorf("totalRead = %d, want 0", totalRead)
		}
		if readingProgress != 0 {
			t.Errorf("readingProgress = %d, want 0", readingProgress)
		}
		if popularAuthor != "N/A" {
			t.Errorf("popularAuthor = %s, want N/A", popularAuthor)
		}
	})

	t.Run("BooksWithDifferentStatuses", func(t *testing.T) {
		_, err := db.ExecContext(ctx, "DELETE FROM books")
		if err != nil {
			t.Fatalf("Failed to clear database: %v", err)
		}

		addBook(uuid.NewString(), "Book 1", "Author A", models.BookComplete)
		addBook(uuid.NewString(), "Book 2", "Author A", models.BookComplete)
		addBook(uuid.NewString(), "Book 3", "Author B", models.BookReading)
		addBook(uuid.NewString(), "Book 4", "Author B", models.BookUnread)

		totalRead, readingProgress, popularAuthor, err := store.GetStats(ctx)
		if err != nil {
			t.Errorf("GetStats failed: %v", err)
		}

		if totalRead != 2 {
			t.Errorf("totalRead = %d, want 2", totalRead)
		}

		if readingProgress != 25 {
			t.Errorf("readingProgress = %d, want 25", readingProgress)
		}

		// NOTE:
		// So the logic is that if there are several authors with equally many *max* books like here where you have 2
		// from Author A and 2 from Author B in your collection. Then it will estbalish the
		// most popular author as the one which name comes first alpabetically
		if popularAuthor != "Author A" {
			t.Errorf("popularAuthor = %s, want Author A", popularAuthor)
		}
	})

	t.Run("SingleAuthor", func(t *testing.T) {

		_, err := db.ExecContext(ctx, "DELETE FROM books")
		if err != nil {
			t.Fatalf("Failed to clear database: %v", err)
		}

		addBook(uuid.NewString(), "Book 1", "Author X", models.BookComplete)
		addBook(uuid.NewString(), "Book 2", "Author X", models.BookReading)

		totalRead, readingProgress, popularAuthor, err := store.GetStats(ctx)
		if err != nil {
			t.Errorf("GetStats failed: %v", err)
		}
		if totalRead != 1 {
			t.Errorf("totalRead = %d, want 1", totalRead)
		}
		if readingProgress != 50 {
			t.Errorf("readingProgress = %d, want 50", readingProgress)
		}
		if popularAuthor != "Author X" {
			t.Errorf("popularAuthor = %s, want Author X", popularAuthor)
		}
	})
}
