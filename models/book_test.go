package models

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestBook_Validate(t *testing.T) {
	tests := []struct {
		name    string
		book    *Book
		wantErr error
	}{
		{
			name: "ValidBookWithUnreadStatus",
			book: &Book{
				ID:     uuid.NewString(),
				Title:  "The Go Programming Language",
				Author: "Alan Donovan",
				Status: BookUnread,
			},
			wantErr: nil,
		},
		{
			name: "ValidBookWithReadingStatus",
			book: &Book{
				ID:     uuid.NewString(),
				Title:  "Sapiens: A Brief History of Humankind",
				Author: "Yuval Noah Harari",
				Status: BookReading,
			},
			wantErr: nil,
		},
		{
			name: "MissingID",
			book: &Book{
				ID:     "", // NOTE: Maybe this can be checked directly with uuid.Parse() so no need to check for both "" and uuid.Parse()? Not a big deal can check later
				Title:  "Learning Go",
				Author: "Jon Bodner",
				Status: BookComplete,
			},
			wantErr: ErrMissingID,
		},
		{
			name: "InvalidID",
			book: &Book{
				ID:     "not-a-uuid",
				Title:  "Concurrency in Go",
				Author: "Katherine Cox-Buday",
				Status: BookReading,
			},
			wantErr: ErrInvalidID,
		},
		{
			name: "MissingTitle",
			book: &Book{
				ID:     uuid.NewString(),
				Title:  "",
				Author: "Bill Bryson",
				Status: BookUnread,
			},
			wantErr: ErrMissingTitle,
		},
		{
			name: "WhitespaceOnlyTitle",
			book: &Book{
				ID:     uuid.NewString(),
				Title:  "   ",
				Author: "Bill Bryson",
				Status: BookUnread,
			},
			wantErr: ErrMissingTitle,
		},
		{
			name: "MissingAuthor",
			book: &Book{
				ID:     uuid.NewString(),
				Title:  "A Short History of Nearly Everything",
				Author: "",
				Status: BookComplete,
			},
			wantErr: ErrMissingAuthor,
		},
		{
			name: "InvalidStatus",
			book: &Book{
				ID:     uuid.NewString(),
				Title:  "Go in Action",
				Author: "William Kennedy",
				Status: "invalid",
			},
			wantErr: ErrInvalidStatus,
		},
		{
			name: "EmptyStatus",
			book: &Book{
				ID:     uuid.NewString(),
				Title:  "Go in Action",
				Author: "William Kennedy",
				Status: "",
			},
			wantErr: ErrEmptyStatus,
		},
		{
			name: "SanitizeInputs",
			book: &Book{
				ID:     uuid.NewString(),
				Title:  "  The Pragmatic Programmer  ",
				Author: "  Andrew Hunt  ",
				Status: "  READING  ",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.book.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if got := strings.TrimSpace(tt.book.Title); got != tt.book.Title {
					t.Errorf("Validate() did not trim Title: got %q, want %q", tt.book.Title, got)
				}
				if got := strings.TrimSpace(tt.book.Author); got != tt.book.Author {
					t.Errorf("Validate() did not trim Author: got %q, want %q", tt.book.Author, got)
				}
				if tt.name == "SanitizeInputs" {
					if tt.book.Title != "The Pragmatic Programmer" {
						t.Errorf("Validate() Title = %q, want %q", tt.book.Title, "The Pragmatic Programmer")
					}
					if tt.book.Author != "Andrew Hunt" {
						t.Errorf("Validate() Author = %q, want %q", tt.book.Author, "Andrew Hunt")
					}
					if tt.book.Status != BookReading {
						t.Errorf("Validate() Status = %q, want %q", tt.book.Status, BookReading)
					}
					t.Logf("Sanitized Title: %q, Author: %q, Status: %q", tt.book.Title, tt.book.Author, tt.book.Status)
				}
			}
			if errors.Is(tt.wantErr, ErrInvalidStatus) {
				expectedMsg := fmt.Sprintf("%v: %s", ErrInvalidStatus, tt.book.Status)
				if err.Error() != expectedMsg {
					t.Errorf("Validate() error message = %q, want %q", err.Error(), expectedMsg)
				}
			}
		})
	}
}

func TestBook_GenerateID(t *testing.T) {
	book := &Book{
		Title:  "The Go Programming Language",
		Author: "Alan Donovan",
		Status: BookUnread,
	}

	err := book.GenerateID()
	if err != nil {
		t.Errorf("GenerateID() error = %v, want no error", err)
	}

	if book.ID == "" {
		t.Errorf("GenerateID() did not set ID, got empty string")
	}
}
