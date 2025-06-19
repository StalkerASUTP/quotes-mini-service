package quote

import (
	"quotes-mini-service/internal/storage"
	"testing"
)

func setupTestDB(t *testing.T) *storage.Db {
	dbPath := ":memory:"
	db, err := storage.NewStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	return db
}

func TestQuotesRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuotesRepository(db)

	author := "Test Author"
	quoteText := "Test Quote"

	quote, err := repo.Save(author, quoteText)
	if err != nil {
		t.Fatalf("failed to save quote: %v", err)
	}

	if quote.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if quote.Author != author {
		t.Errorf("expected author %s, got %s", author, quote.Author)
	}
	if quote.Quote != quoteText {
		t.Errorf("expected quote %s, got %s", quoteText, quote.Quote)
	}
	if quote.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestQuotesRepository_Save_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuotesRepository(db)

	author := "Test Author"
	quoteText := "Test Quote"

	// First save
	_, err := repo.Save(author, quoteText)
	if err != nil {
		t.Fatalf("failed to save quote: %v", err)
	}

	// Second save (duplicate)
	_, err = repo.Save(author, quoteText)
	if err == nil {
		t.Error("expected error for duplicate entry")
	}
}

func TestQuotesRepository_GetAllParam(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuotesRepository(db)

	// Add test data
	testQuotes := []struct {
		author string
		quote  string
	}{
		{"Author1", "Quote1"},
		{"Author2", "Quote2"},
		{"Author1", "Quote3"},
	}

	for _, tq := range testQuotes {
		_, err := repo.Save(tq.author, tq.quote)
		if err != nil {
			t.Fatalf("failed to save test quote: %v", err)
		}
	}
	quotes, count, err := repo.GetAllParam("")
	if err != nil {
		t.Fatalf("failed to get all quotes: %v", err)
	}
	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
	if len(quotes) != 3 {
		t.Errorf("expected 3 quotes, got %d", len(quotes))
	}

	quotes, count, err = repo.GetAllParam("Author1")
	if err != nil {
		t.Fatalf("failed to get quotes by author: %v", err)
	}
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
	if len(quotes) != 2 {
		t.Errorf("expected 2 quotes, got %d", len(quotes))
	}
}

func TestQuotesRepository_GetRandom(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuotesRepository(db)

	newQuote, err := repo.Save("Test Author", "Test Quote")
	if err != nil {
		t.Fatalf("failed to save test quote: %v", err)
	}

	quote, err := repo.GetRandom()
	if err != nil {
		t.Fatalf("failed to get random quote: %v", err)
	}

	if newQuote.ID != quote.ID {
		t.Error("expected non-zero ID")
	}
	if quote.Author != newQuote.Author {
		t.Error("expected non-empty author")
	}
	if quote.Quote != newQuote.Quote {
		t.Error("expected non-empty quote")
	}
}

func TestQuotesRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuotesRepository(db)

	quote, err := repo.Save("Test Author", "Test Quote")
	if err != nil {
		t.Fatalf("failed to save test quote: %v", err)
	}

	err = repo.Delete(quote.ID)
	if err != nil {
		t.Fatalf("failed to delete quote: %v", err)
	}

	_, count, err := repo.GetAllParam("")
	if err != nil {
		t.Fatalf("failed to get quotes: %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
}

func TestQuotesRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuotesRepository(db)

	err := repo.Delete(999)
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}
