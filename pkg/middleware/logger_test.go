package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	rw := NewResponseWriter(w)

	if rw.StatusCode() != 200 {
		t.Errorf("expected initial status 200, got %d", rw.StatusCode())
	}
	if rw.BytesWritten() != 0 {
		t.Errorf("expected initial bytes written 0, got %d", rw.BytesWritten())
	}
	rw.WriteHeader(404)
	if rw.StatusCode() != 404 {
		t.Errorf("expected status 404, got %d", rw.StatusCode())
	}
	data := []byte("test data")
	n, err := rw.Write(data)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != len(data) {
		t.Errorf("expected %d bytes written, got %d", len(data), n)
	}
	if rw.BytesWritten() != len(data) {
		t.Errorf("expected bytes written %d, got %d", len(data), rw.BytesWritten())
	}
}

func TestMiddleware(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	middleware := New(log)
	wrappedHandler := middleware(handler)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if body != "test response" {
		t.Errorf("expected body 'test response', got '%s'", body)
	}
}