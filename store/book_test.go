package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"book-tracker/models"

	"github.com/google/uuid"
)

func setupDB(t *testing.T) (*sql.DB, func()) {
	db, cleanup, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	return db, cleanup
}

func TestBookStore(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	store := NewBookStore(db)
	ctx := context.Background()

	t.Run("CreateBook", func(t *testing.T) {
		book := &models.Book{
			ID:     uuid.NewString(),
			Title:  "Test Book",
			Author: "Test Author",
			Status: models.BookUnread,
		}
		err := store.CreateBook(ctx, book)
		if err != nil {
			t.Fatalf("CreateBook failed: %v", err)
		}

		got, err := store.GetBook(ctx, book.ID)
		if err != nil {
			t.Fatalf("GetBook failed: %v", err)
		}
		if got.ID != book.ID || got.Title != book.Title || got.Author != book.Author || got.Status != book.Status {
			t.Errorf("GetBook = %+v, want %+v", got, book)
		}
	})

	t.Run("GetBook", func(t *testing.T) {
		book := &models.Book{
			ID:     uuid.NewString(),
			Title:  "Another Book",
			Author: "Another Author",
			Status: models.BookReading,
		}
		err := store.CreateBook(ctx, book)
		if err != nil {
			t.Fatalf("CreateBook failed: %v", err)
		}

		got, err := store.GetBook(ctx, book.ID)
		if err != nil {
			t.Errorf("GetBook failed: %v", err)
		}
		if got.ID != book.ID || got.Title != book.Title || got.Author != book.Author || got.Status != book.Status {
			t.Errorf("GetBook = %+v, want %+v", got, book)
		}

		_, err = store.GetBook(ctx, "nonexistent")
		if !errors.Is(err, ErrBookNotFound) {
			t.Errorf("GetBook error = %v, want %v", err, ErrBookNotFound)
		}
	})

	t.Run("ListBooks", func(t *testing.T) {
		_, err := db.ExecContext(ctx, "DELETE FROM books")
		if err != nil {
			t.Fatalf("Failed to clear database: %v", err)
		}

		books := []models.Book{
			{ID: uuid.NewString(), Title: "Book A", Author: "Author A", Status: models.BookUnread},
			{ID: uuid.NewString(), Title: "Book B", Author: "Author B", Status: models.BookReading},
		}
		for _, b := range books {
			err := store.CreateBook(ctx, &b)
			if err != nil {
				t.Fatalf("CreateBook failed: %v", err)
			}
		}

		got, err := store.ListBooks(ctx, "", 10, 0, "", "")
		if err != nil {
			t.Errorf("ListBooks failed: %v", err)
		}
		if len(got) != 2 {
			t.Errorf("ListBooks returned %d books, want 2", len(got))
		}
		if got[0].Title != "Book A" || got[1].Title != "Book B" {
			t.Errorf("ListBooks titles = %v, %v, want Book A, Book B", got[0].Title, got[1].Title)
		}
	})

	t.Run("UpdateBook", func(t *testing.T) {
		book := &models.Book{
			ID:     uuid.NewString(),
			Title:  "Original Book",
			Author: "Original Author",
			Status: models.BookUnread,
		}
		err := store.CreateBook(ctx, book)
		if err != nil {
			t.Fatalf("CreateBook failed: %v", err)
		}

		updatedBook := &models.Book{
			ID:     book.ID,
			Title:  "Updated Book",
			Author: "Updated Author",
			Status: models.BookReading,
		}
		err = store.UpdateBook(ctx, updatedBook)
		if err != nil {
			t.Errorf("UpdateBook failed: %v", err)
		}
		got, err := store.GetBook(ctx, book.ID)
		if err != nil {
			t.Fatalf("GetBook failed: %v", err)
		}
		if got.ID != updatedBook.ID || got.Title != updatedBook.Title || got.Author != updatedBook.Author || got.Status != updatedBook.Status {
			t.Errorf("GetBook after update = %+v, want %+v", got, updatedBook)
		}

		nonExistentBook := &models.Book{
			ID:     "nonexistent",
			Title:  "No Book",
			Author: "No Author",
			Status: models.BookUnread,
		}
		err = store.UpdateBook(ctx, nonExistentBook)
		if !errors.Is(err, ErrBookNotFound) {
			t.Errorf("UpdateBook error = %v, want %v", err, ErrBookNotFound)
		}
	})

	t.Run("DeleteBook", func(t *testing.T) {
		book := &models.Book{
			ID:     uuid.NewString(),
			Title:  "Delete Book",
			Author: "Delete Author",
			Status: models.BookUnread,
		}
		err := store.CreateBook(ctx, book)
		if err != nil {
			t.Fatalf("CreateBook failed: %v", err)
		}

		err = store.DeleteBook(ctx, book.ID)
		if err != nil {
			t.Errorf("DeleteBook failed: %v", err)
		}
		_, err = store.GetBook(ctx, book.ID)
		if !errors.Is(err, ErrBookNotFound) {
			t.Errorf("GetBook after delete should return ErrBookNotFound, got %v", err)
		}
	})

	t.Run("CountBooks", func(t *testing.T) {
		_, err := db.ExecContext(ctx, "DELETE FROM books")
		if err != nil {
			t.Fatalf("Failed to clear database: %v", err)
		}

		books := []models.Book{
			{ID: uuid.NewString(), Title: "Count Book", Author: "Count Author", Status: models.BookUnread},
			{ID: uuid.NewString(), Title: "Count Book", Author: "Count Author", Status: models.BookReading},
		}
		for _, b := range books {
			err := store.CreateBook(ctx, &b)
			if err != nil {
				t.Fatalf("CreateBook failed: %v", err)
			}
		}

		total, byStatus, err := store.CountBooks(ctx)
		if err != nil {
			t.Errorf("CountBooks failed: %v", err)
		}
		if total != 2 {
			t.Errorf("CountBooks total = %d, want 2", total)
		}
		if byStatus[string(models.BookUnread)] != 1 || byStatus[string(models.BookReading)] != 1 {
			t.Errorf("CountBooks byStatus = %+v, want unread:1, reading:1", byStatus)
		}
	})

	t.Run("ListBooks_Pagination", func(t *testing.T) {
		_, err := db.ExecContext(ctx, "DELETE FROM books")
		if err != nil {
			t.Fatalf("Failed to clear database: %v", err)
		}

		books := []models.Book{
			{ID: uuid.NewString(), Title: "Book 1", Author: "Author 1", Status: models.BookUnread},
			{ID: uuid.NewString(), Title: "Book 2", Author: "Author 2", Status: models.BookReading},
			{ID: uuid.NewString(), Title: "Book 3", Author: "Author 3", Status: models.BookUnread},
			{ID: uuid.NewString(), Title: "Book 4", Author: "Author 4", Status: models.BookReading},
			{ID: uuid.NewString(), Title: "Book 5", Author: "Author 5", Status: models.BookUnread},
		}
		for _, b := range books {
			err := store.CreateBook(ctx, &b)
			if err != nil {
				t.Fatalf("CreateBook failed: %v", err)
			}
		}

		tests := []struct {
			name       string
			limit      int
			offset     int
			wantCount  int
			wantTitles []string
		}{
			{
				name:       "FirstPage_Limit3",
				limit:      3,
				offset:     0,
				wantCount:  3,
				wantTitles: []string{"Book 1", "Book 2", "Book 3"},
			},
			{
				name:       "SecondPage_Limit3",
				limit:      3,
				offset:     3,
				wantCount:  2,
				wantTitles: []string{"Book 4", "Book 5"},
			},
			{
				name:       "BeyondTotal_Limit3",
				limit:      3,
				offset:     6,
				wantCount:  0,
				wantTitles: []string{},
			},
			{
				name:       "LimitZero",
				limit:      0,
				offset:     0,
				wantCount:  0,
				wantTitles: []string{},
			},
			{
				name:       "FullList_Limit10",
				limit:      10,
				offset:     0,
				wantCount:  5,
				wantTitles: []string{"Book 1", "Book 2", "Book 3", "Book 4", "Book 5"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := store.ListBooks(ctx, "", tt.limit, tt.offset, "", "")
				if err != nil {
					t.Errorf("ListBooks failed: %v", err)
				}
				if len(got) != tt.wantCount {
					t.Errorf("ListBooks returned %d books, want %d", len(got), tt.wantCount)
				}
				if len(got) != len(tt.wantTitles) {
					t.Errorf("ListBooks returned %d books, but wantTitles has %d", len(got), len(tt.wantTitles))
				}
				for i, title := range tt.wantTitles {
					if i >= len(got) {
						break
					}
					if got[i].Title != title {
						t.Errorf("ListBooks title %d = %s, want %s", i, got[i].Title, title)
					}
				}
			})
		}
	})
}
