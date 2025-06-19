package quote

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"
	"quotes-mini-service/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type QuotesRepository struct {
	Database *storage.Db
}

func NewQuotesRepository(databse *storage.Db) *QuotesRepository {
	return &QuotesRepository{
		Database: databse,
	}
}
func (repo *QuotesRepository) Save(authorSave, quoteSave string) (*Quote, error) {
	const op = "quote.repository.Create"
	tx, err := repo.Database.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback()
	insertStmt, err := tx.Prepare("INSERT INTO quotes(author,quote) VALUES(?, ?)")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer insertStmt.Close()
	res, err := insertStmt.Exec(authorSave, quoteSave)
	if err != nil {
		if isDuplicateError(err) {
			return nil, fmt.Errorf("%s: duplicate entry: %w", op, err)
		}
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	var quotes Quote
	selectStmt, err := tx.Prepare("SELECT * FROM quotes WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare select: %w", op, err)
	}
	defer selectStmt.Close()
	err = selectStmt.QueryRow(id).
		Scan(&quotes.ID,
			&quotes.Author,
			&quotes.Quote,
			&quotes.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: scan result: %w", op, err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: commit transaction: %w", op, err)
	}
	return &quotes, nil
}

func (repo *QuotesRepository) GetAllParam(author string) ([]Quote, int, error) {
	const op = "quote.repository.GetAllParan"
	var (
		rows *sql.Rows
		err  error
	)
	tx, err := repo.Database.Begin()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback()
	if author == "" {
		rows, err = tx.Query("SELECT * FROM quotes")
	} else {
		rows, err = tx.Query("SELECT * FROM quotes WHERE author = ?", author)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query execution: %w", op, err)
	}
	defer rows.Close()
	var quotes []Quote
	for rows.Next() {
		var quote Quote
		if err := rows.Scan(&quote.ID,
			&quote.Author,
			&quote.Quote,
			&quote.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("%s: commit transaction: %w", op, err)
		}
		quotes = append(quotes, quote)
	}
	if err = tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("%s: commit transaction: %w", op, err)
	}
	return quotes, len(quotes), nil
}

func (repo *QuotesRepository) GetRandom() (*Quote, error) {
	const op = "quote.repository.GetRandom"
	tx, err := repo.Database.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback()
	var count int
	err = tx.QueryRow("SELECT count_value FROM counters WHERE table_name = 'quotes'").
		Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to get count: %w", err)
	}
	offset := rand.IntN(count)
	var randomQuote Quote
	err = tx.QueryRow("SELECT * FROM quotes ORDER BY id LIMIT 1 OFFSET ?", offset).
		Scan(&randomQuote.ID,
			&randomQuote.Author,
			&randomQuote.Quote,
			&randomQuote.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: scan result: %w", op, err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: commit transaction: %w", op, err)
	}
	return &randomQuote, nil
}

func (repo *QuotesRepository) Delete(id int) error {
	const op = "quote.repository.Delete"
	tx, err := repo.Database.Begin()
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback()
	result, err := tx.Exec("DELETE FROM quotes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("%s: delete operation: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows affection: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: quote with id %d not found", op, id)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}
	return nil
}

func isDuplicateError(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	return false
}
