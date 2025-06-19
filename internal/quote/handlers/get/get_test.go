package get

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"quotes-mini-service/internal/quote"
	"testing"
	"time"
)

type mockQuoteGetter struct {
	quotes []quote.Quote
}

func newMockQuoteGetter() *mockQuoteGetter {
	return &mockQuoteGetter{
		quotes: []quote.Quote{
			{ID: 1, Author: "Author1", Quote: "Quote1", CreatedAt: time.Now()},
			{ID: 2, Author: "Author2", Quote: "Quote2", CreatedAt: time.Now()},
			{ID: 3, Author: "Author1", Quote: "Quote3", CreatedAt: time.Now()},
		},
	}
}

func (m *mockQuoteGetter) GetAllParam(author string) ([]quote.Quote, int, error) {
	var filtered []quote.Quote
	for _, q := range m.quotes {
		if author == "" || q.Author == author {
			filtered = append(filtered, q)
		}
	}
	return filtered, len(filtered), nil
}

func (m *mockQuoteGetter) GetRandom() (*quote.Quote, error) {
	if len(m.quotes) == 0 {
		return nil, errors.New("no quotes available")
	} 
	return &m.quotes[0], nil
}

func TestAllParam_Success(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	getter := newMockQuoteGetter()
	handler := AllParam(log, getter)

	r := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response GetWithParamResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Quotes) != 3 {
		t.Errorf("expected 3 quotes, got %d", len(response.Quotes))
	}
	if response.Count != 3 {
		t.Errorf("expected count 3, got %d", response.Count)
	}
}

func TestAllParam_WithAuthorFilter(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	getter := newMockQuoteGetter()
	handler := AllParam(log, getter)

	r := httptest.NewRequest(http.MethodGet, "/quotes?author=Author1", nil)
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response GetWithParamResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Quotes) != 2 {
		t.Errorf("expected 2 quotes, got %d", len(response.Quotes))
	}
	if response.Count != 2 {
		t.Errorf("expected count 2, got %d", response.Count)
	}
}

func TestRandom_Success(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	getter := newMockQuoteGetter()
	handler := Random(log, getter)

	r := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response quote.Quote
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID == 0 {
		t.Error("expected non-zero ID")
	}
}