package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"quotes-mini-service/internal/quote"
	"strings"
	"testing"
	"time"
)

type mockQuoteSaver struct {
	quotes map[string]*quote.Quote
	lastId int
}

func newMockQuoteSaver() *mockQuoteSaver {
	return &mockQuoteSaver{
		quotes: make(map[string]*quote.Quote),
		lastId: 0,
	}
}
func (m *mockQuoteSaver) Save(author, quoteText string) (*quote.Quote, error) {
	key := author + ":" + quoteText
	if _, exists := m.quotes[key]; exists {
		return nil, errors.New("duplicate entry")
	}
	newid := m.lastId + 1
	q := &quote.Quote{
		ID:        newid,
		Author:    author,
		Quote:     quoteText,
		CreatedAt: time.Now(),
	}
	m.lastId = newid
	m.quotes[key] = q
	return q, nil
}

func TestSave_Success(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := newMockQuoteSaver()
	handler := New(log, repo)

	req := Request{
		Author: "Test Author",
		Quote:  "Test Quote",
	}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response quote.Quote
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Author != req.Author {
		t.Errorf("expected author %s, got %s", req.Author, response.Author)
	}
	if response.Quote != req.Quote {
		t.Errorf("expected quote %s, got %s", req.Quote, response.Quote)
	}
}

func TestSave_InvalidJSON(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := newMockQuoteSaver()
	handler := New(log, repo)

	r := httptest.NewRequest(http.MethodPost, "/quotes", strings.NewReader("invalid json"))
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSave_MissingAuthor(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := newMockQuoteSaver()
	handler := New(log, repo)

	req := Request{
		Quote: "Test Quote",
	}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSave_DuplicateEntry(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := newMockQuoteSaver()
	handler := New(log, repo)

	req := Request{
		Author: "Test Author",
		Quote:  "Test Quote",
	}
	body, _ := json.Marshal(req)

	r1 := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	w1 := httptest.NewRecorder()
	handler(w1, r1)

	r2 := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	w2 := httptest.NewRecorder()
	handler(w2, r2)

	if w2.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w2.Code)
	}
}
