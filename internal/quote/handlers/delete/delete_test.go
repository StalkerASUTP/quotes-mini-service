package delete

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type mockQuoteDeleter struct {
	existingIDs map[int]bool
}

func newMockQuoteDeleter() *mockQuoteDeleter {
	return &mockQuoteDeleter{
		existingIDs: map[int]bool{
			1: true,
			2: true,
			3: true,
		},
	}
}

func (m *mockQuoteDeleter) Delete(id int) error {
	if !m.existingIDs[id] {
		return errors.New("not found")
	}
	delete(m.existingIDs, id)
	return nil
}

func TestDelete_Success(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	deleter := newMockQuoteDeleter()
	handler := New(log, deleter)

	r := httptest.NewRequest(http.MethodDelete, "/quotes/1", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestDelete_InvalidID(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	deleter := newMockQuoteDeleter()
	handler := New(log, deleter)

	r := httptest.NewRequest(http.MethodDelete, "/quotes/invalid", nil)
	r.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDelete_NotFound(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	deleter := newMockQuoteDeleter()
	handler := New(log, deleter)

	r := httptest.NewRequest(http.MethodDelete, "/quotes/999", nil)
	r.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}