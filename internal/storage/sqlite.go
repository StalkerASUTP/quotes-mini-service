package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrAuthorNotFound = errors.New("author not found")
	ErrIdNotFound        = errors.New("id not found")
)

type Db struct {
	*sql.DB
}

func NewStorage(dbPath string) (*Db, error) {
	const op = "storage.NewStorage"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	createTables := []string{
		`CREATE TABLE IF NOT EXISTS quotes(
		id INTEGER PRIMARY KEY,
		author TEXT NOT NULL,
		quote TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(author, quote)
	)`,
		` CREATE TABLE IF NOT EXISTS counters(
            table_name TEXT PRIMARY KEY,
            count_value INTEGER DEFAULT 0
    )`,
		`CREATE INDEX IF NOT EXISTS idx_author ON quotes(author)`,
		`INSERT OR IGNORE INTO counters (table_name, count_value) 
         VALUES ('quotes', 0)`,
	}
	for _, query := range createTables {
		if _, err := db.Exec(query); err != nil {
			db.Close()
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	createTrigger := []string{
		`CREATE TRIGGER IF NOT EXISTS update_quotes_counter
    AFTER INSERT ON quotes
    FOR EACH ROW
    BEGIN
        UPDATE counters SET count_value = count_value + 1 
        WHERE table_name = 'quotes';
    END;`,
		`CREATE TRIGGER IF NOT EXISTS delete_quotes_counter
    AFTER DELETE ON quotes
    FOR EACH ROW
    BEGIN
        UPDATE counters SET count_value = count_value - 1 
        WHERE table_name = 'quotes';
    END;
		`,
	}
	for _, query := range createTrigger {
		if _, err := db.Exec(query); err != nil {
			db.Close()
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return &Db{db}, nil
}
