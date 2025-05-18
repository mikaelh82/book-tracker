package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const createBooksTable = `
CREATE TABLE IF NOT EXISTS books (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    status TEXT NOT NULL
)
`

// Documentation: https://duckdb.org/docs/stable/sql/statements/create_index.html
// Just to speed up like GET /books?status=reading for the major book worm!!! =)
const createStatusIndex = `
CREATE INDEX IF NOT EXISTS idx_status ON books (status)
`

func NewDB(dbPath string) (*sql.DB, func(), error) {
	// Documentation: https://duckdb.org/docs/stable/clients/go.html
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, nil, fmt.Errorf("open duckdb: %w", err)
	}

	_, err = db.Exec(createBooksTable)

	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("create table: %w", err)
	}

	_, err = db.Exec(createStatusIndex)

	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("create index: %w", err)
	}

	return db, func() { db.Close() }, nil // return db.Close() in a wrapper to defer closing the db conenction
}
