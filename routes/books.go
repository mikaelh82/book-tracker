package routes

import (
	"book-tracker/handlers"
	"net/http"
	"strings"
)

func SetupBooksRoutes(mux *http.ServeMux, handler *handlers.BookHandler) {
	// NOTE:
	// Handle POST and GET /api/v1/books.
	mux.HandleFunc("/api/v1/books", func(w http.ResponseWriter, r *http.Request) { // NOTE: this version can later be linked to config and not hardcoded
		switch r.Method {
		case "POST":
			handler.CreateBook(w, r)
		case "GET":
			handler.ListBooks(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// NOTE:
	// Handle PUT and DELETE /api/v1/books/{id}.
	mux.HandleFunc("/api/v1/books/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if !strings.HasPrefix(path, "/api/v1/books/") { // NOTE: this version can later be linked to config and not hardcoded
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		id := strings.TrimPrefix(path, "/api/v1/books/")
		if id == "" {
			http.Error(w, "Book ID required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "PUT":
			handler.UpdateBook(w, r, id)
		case "DELETE":
			handler.DeleteBook(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
