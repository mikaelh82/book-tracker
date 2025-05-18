package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrMissingID     = errors.New("id is mising")
	ErrInvalidID     = errors.New("id is invalid")
	ErrMissingTitle  = errors.New("title is missing")
	ErrMissingAuthor = errors.New("author is missing")
	ErrInvalidStatus = errors.New("invalid status: must be unread, reading or complete")
	ErrEmptyStatus   = errors.New("status cannot be empty")
)

type BookStatus string

const (
	BookUnread   BookStatus = "unread"
	BookReading  BookStatus = "reading"
	BookComplete BookStatus = "complete"
)

type Book struct {
	ID     string     `json:"id"`
	Title  string     `json:"title"`
	Author string     `json:"author"`
	Status BookStatus `json:"status"`
}

func (b *Book) Validate() error {
	b.Title = strings.TrimSpace(b.Title)
	b.Author = strings.TrimSpace(b.Author)

	if b.ID == "" {
		return ErrMissingID
	}

	if b.ID != "" {
		if _, err := uuid.Parse(b.ID); err != nil {
			return ErrInvalidID
		}
	}

	if b.Title == "" {
		return ErrMissingTitle
	}

	if b.Author == "" {
		return ErrMissingAuthor
	}

	trimmedStatus := strings.TrimSpace(string(b.Status))
	if trimmedStatus == "" {
		return ErrEmptyStatus
	}

	status := BookStatus(strings.ToLower(trimmedStatus))

	switch status {
	case BookUnread, BookReading, BookComplete:
		// ALL GOOD
		b.Status = status
	default:
		return fmt.Errorf("%w: %s", ErrInvalidStatus, status) // NOTE: %w is a Go feature. Wrapping an existing error
	}

	return nil
}

func (b *Book) GenerateID() error {
	id, err := uuid.NewRandom()

	if err != nil {
		return fmt.Errorf("failed to generate id: %w", err)
	}

	b.ID = id.String()

	return nil
}
