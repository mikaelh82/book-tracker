package store

import (
	"database/sql"
	"fmt"

	_ "github.com/marcboeker/go-duckdb"
)

const createBooksTable = `
CREATE TABLE IF NOT EXIST books (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	author TEXT NOT NULL,
	STATUS TEXT NOT NULL,
	UNIQUE(title, author)
)
`

// Documentation: https://duckdb.org/docs/stable/sql/statements/create_index.html
// Just to speed up like GET /books?status=reading for the major book worm!!! =)
const createStatusIndex = `
CREATE INDEX IF NOT EXISTS idx_status ON books (status)
`

func NewDB(dbPath string) (*sql.DB, func(), error) {
	// Documentation: https://duckdb.org/docs/stable/clients/go.html
	db, err := sql.Open("duckdb", dbPath)

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
