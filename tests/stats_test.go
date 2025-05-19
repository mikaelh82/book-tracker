package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"book-tracker/handlers"
	"book-tracker/models"
	"book-tracker/routes"
	"book-tracker/services"
	"book-tracker/store"
)

func setupStats(t *testing.T) (*http.ServeMux, store.BookStore, func()) {
	db, closeDB, err := store.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize SQLite: %v", err)
	}
	statsStore := store.NewStatsStore(db)
	bookStore := store.NewBookStore(db)
	statsService := services.NewStatsService(statsStore)
	statsHandler := handlers.NewStatsHandler(statsService)
	mux := http.NewServeMux()
	routes.SetupStatsRoutes(mux, statsHandler)
	return mux, bookStore, closeDB
}

func TestStatsRoutes(t *testing.T) {
	t.Run("GET_EmptyDatabase", func(t *testing.T) {
		mux, _, closeDB := setupStats(t)
		defer closeDB()

		req, _ := http.NewRequest("GET", "/api/v1/stats", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var stats models.Stats
		if err := json.NewDecoder(rr.Body).Decode(&stats); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		if stats.TotalRead != 0 || stats.ReadingProgress != 0 || stats.PopularAuthor != "N/A" {
			t.Errorf("Expected empty stats {0, 0, 'N/A'}, got %+v", stats)
		}
	})

	t.Run("GET_WithBooks", func(t *testing.T) {
		mux, bookStore, closeDB := setupStats(t)
		defer closeDB()

		books := []models.Book{
			{Title: "Book A", Author: "Jane Austen", Status: models.BookComplete},
			{Title: "Book B", Author: "Jane Austen", Status: models.BookReading},
			{Title: "Book C", Author: "Mark Twain", Status: models.BookUnread},
		}
		for i := range books {
			if err := books[i].GenerateID(); err != nil {
				t.Fatalf("Failed to generate UUID for book %d: %v", i, err)
			}
			if err := books[i].Validate(); err != nil {
				t.Fatalf("Invalid test book %d: %v", i, err)
			}
			if err := bookStore.CreateBook(context.Background(), &books[i]); err != nil {
				t.Fatalf("Failed to seed book %d: %v", i, err)
			}
		}

		req, _ := http.NewRequest("GET", "/api/v1/stats", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var stats models.Stats
		if err := json.NewDecoder(rr.Body).Decode(&stats); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		if stats.TotalRead != 1 {
			t.Errorf("Expected total_read 1, got %d", stats.TotalRead)
		}
		if stats.ReadingProgress != 33 { // 1 reading / 3 total = 33%
			t.Errorf("Expected reading_progress 33, got %d", stats.ReadingProgress)
		}
		if stats.PopularAuthor != "Jane Austen" {
			t.Errorf("Expected popular_author 'Jane Austen', got %s", stats.PopularAuthor)
		}
	})

	t.Run("InvalidMethod_Stats", func(t *testing.T) {
		mux, _, closeDB := setupStats(t)
		defer closeDB()

		req, _ := http.NewRequest("POST", "/api/v1/stats", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", rr.Code)
		}
	})
}
