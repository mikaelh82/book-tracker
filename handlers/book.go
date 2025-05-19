package handlers

import (
	"book-tracker/models"
	"book-tracker/services"
	"book-tracker/store"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type BookHandler struct {
	service services.BookService
}

func NewBookHandler(service services.BookService) *BookHandler {
	return &BookHandler{service: service}
}

type errorResponse struct {
	Error string `json:"error"`
}

// Just a small utility for error serialization/encoding
func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	// bug here. Interesting. For Go to know how to decode a BookStatus (completed etc). it needs to implement a UnmarshalJSON data
	// Documentation: https://dev.to/arshamalh/how-to-unmarshal-json-in-a-custom-way-in-golang-42m5
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request: %v", err))
		return
	}
	if err := h.service.CreateBook(r.Context(), &book); err != nil {
		if errors.Is(err, store.ErrBookNotFound) {
			writeError(w, http.StatusConflict, err)
		} else if errors.Is(err, models.ErrMissingTitle) || errors.Is(err, models.ErrMissingAuthor) || errors.Is(err, models.ErrInvalidStatus) {
			writeError(w, http.StatusBadRequest, err)
		} else {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("create book error: %v", err))
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(book); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to encode response"))
	}
}

func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	// here, we could also later implement author and title if we wish
	// the store has the functionality for it but its whitespaced in the service

	limit := 10
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 { // this is cardcoded for now but this can later be moved to Config
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid limit: must be a number between 1 and 1000"))
			return
		}
	}
	offset := 0
	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid offset: must be a non-negative number"))
			return
		}
	}

	books, err := h.service.ListBooks(r.Context(), status, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("list books error: %v", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(books); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to encode response"))
	}
}
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request, id string) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request: %v", err))
		return
	}
	book.ID = id
	if err := h.service.UpdateBook(r.Context(), &book); err != nil {
		if errors.Is(err, store.ErrBookNotFound) {
			writeError(w, http.StatusNotFound, err)
		} else if errors.Is(err, models.ErrMissingTitle) || errors.Is(err, models.ErrMissingAuthor) || errors.Is(err, models.ErrInvalidStatus) {
			writeError(w, http.StatusBadRequest, err)
		} else {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("update book error: %v", err))
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to encode response"))
	}
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteBook(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrBookNotFound) {
			writeError(w, http.StatusNotFound, err)
		} else {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("delete book error: %v", err))
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
