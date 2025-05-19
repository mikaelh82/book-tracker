package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"book-tracker/handlers"
	"book-tracker/models"
	"book-tracker/routes"
	"book-tracker/services"
	"book-tracker/store"
)

func setupBooks(t *testing.T) (*http.ServeMux, store.BookStore, func()) {
	db, closeDB, err := store.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize SQLite: %v", err)
	}
	bookStore := store.NewBookStore(db)
	bookService := services.NewBookService(bookStore)
	bookHandler := handlers.NewBookHandler(bookService)
	mux := http.NewServeMux()
	routes.SetupBooksRoutes(mux, bookHandler)
	return mux, bookStore, closeDB
}

func TestBooksRoutes(t *testing.T) {
	t.Run("POST_CreateBook", func(t *testing.T) {
		mux, bookStore, closeDB := setupBooks(t)
		defer closeDB()

		book := models.Book{
			Title:  "Test Book",
			Author: "Test Author",
			Status: models.BookUnread,
		}
		body, _ := json.Marshal(book)
		req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", rr.Code)
		}

		var createdBook models.Book
		if err := json.NewDecoder(rr.Body).Decode(&createdBook); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		if createdBook.Title != book.Title || createdBook.Author != book.Author || createdBook.Status != book.Status {
			t.Errorf("Response book = %+v, want %+v", createdBook, book)
		}

		dbBook, err := bookStore.GetBook(context.Background(), createdBook.ID)
		if err != nil {
			t.Errorf("Failed to get book from store: %v", err)
		}
		if dbBook.Title != book.Title {
			t.Errorf("Stored book title = %s, want %s", dbBook.Title, book.Title)
		}
	})

	t.Run("GET_ListBooks", func(t *testing.T) {
		mux, bookStore, closeDB := setupBooks(t)
		defer closeDB()

		books := []models.Book{
			{Title: "Book A", Author: "Author A", Status: models.BookUnread},
			{Title: "Book B", Author: "Author B", Status: models.BookReading},
		}
		for i := range books {
			if err := books[i].GenerateID(); err != nil {
				t.Fatalf("Failed to generate UUID: %v", err)
			}
			if err := bookStore.CreateBook(context.Background(), &books[i]); err != nil {
				t.Fatalf("Failed to seed book: %v", err)
			}
		}

		req, _ := http.NewRequest("GET", "/api/v1/books?limit=10&offset=0", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var listedBooks []models.Book
		if err := json.NewDecoder(rr.Body).Decode(&listedBooks); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		if len(listedBooks) != 2 {
			t.Errorf("Expected 2 books, got %d", len(listedBooks))
		}
		if listedBooks[0].Title != "Book A" || listedBooks[1].Title != "Book B" {
			t.Errorf("Titles = %s, %s; want Book A, Book B", listedBooks[0].Title, listedBooks[1].Title)
		}
	})

	t.Run("PUT_UpdateBook", func(t *testing.T) {
		mux, bookStore, closeDB := setupBooks(t)
		defer closeDB()

		book := models.Book{
			Title:  "Original Book",
			Author: "Original Author",
			Status: models.BookUnread,
		}
		if err := book.GenerateID(); err != nil {
			t.Fatalf("Failed to generate UUID: %v", err)
		}
		if err := bookStore.CreateBook(context.Background(), &book); err != nil {
			t.Fatalf("Failed to seed book: %v", err)
		}

		updatedBook := models.Book{
			Title:  "Updated Book",
			Author: "Updated Author",
			Status: models.BookReading,
		}
		body, _ := json.Marshal(updatedBook)
		req, _ := http.NewRequest("PUT", "/api/v1/books/"+book.ID, bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var responseBook models.Book
		if err := json.NewDecoder(rr.Body).Decode(&responseBook); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		if responseBook.Title != updatedBook.Title || responseBook.Status != updatedBook.Status {
			t.Errorf("Response book = %+v, want %+v", responseBook, updatedBook)
		}
	})

	t.Run("DELETE_DeleteBook", func(t *testing.T) {
		mux, bookStore, closeDB := setupBooks(t)
		defer closeDB()

		// Seed database
		book := models.Book{
			Title:  "Delete Book",
			Author: "Delete Author",
			Status: models.BookUnread,
		}
		if err := book.GenerateID(); err != nil {
			t.Fatalf("Failed to generate UUID: %v", err)
		}
		if err := bookStore.CreateBook(context.Background(), &book); err != nil {
			t.Fatalf("Failed to seed book: %v", err)
		}

		req, _ := http.NewRequest("DELETE", "/api/v1/books/"+book.ID, nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", rr.Code)
		}

		_, err := bookStore.GetBook(context.Background(), book.ID)
		if !errors.Is(err, store.ErrBookNotFound) {
			t.Errorf("Expected ErrBookNotFound, got %v", err)
		}
	})

	t.Run("InvalidMethod_Books", func(t *testing.T) {
		mux, _, closeDB := setupBooks(t)
		defer closeDB()

		req, _ := http.NewRequest("PATCH", "/api/v1/books", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", rr.Code)
		}
	})

	t.Run("InvalidPath_BookID", func(t *testing.T) {
		mux, _, closeDB := setupBooks(t)
		defer closeDB()

		req, _ := http.NewRequest("PUT", "/api/v1/books/", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})
}
