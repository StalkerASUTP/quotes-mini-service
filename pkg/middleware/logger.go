package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type ResponseWriter interface {
	http.ResponseWriter
	StatusCode() int
	BytesWritten() int
}
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
	wroteHeader  bool
}

func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     200,
	}
}
func (w *responseWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.statusCode = code
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
}
func (w *responseWriter) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(200)
	}
	n, err := w.ResponseWriter.Write(data)
	w.bytesWritten += n
	return n, err
}
func (w *responseWriter) StatusCode() int {
	return w.statusCode
}
func (w *responseWriter) BytesWritten() int {
	return w.bytesWritten
}

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)
		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)
			wr := NewResponseWriter(w)
			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", wr.StatusCode()),
					slog.Int("bytes", wr.BytesWritten()),
					slog.String("duration", time.Since(t1).String()))
			}()
			next.ServeHTTP(wr, r)
		}
		return http.HandlerFunc(fn)
	}
}
