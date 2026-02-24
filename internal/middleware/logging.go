package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type contextKey string
		const requestIDKey contextKey = "requestID"

		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", requestID)

		start := time.Now()

		wrapped := &statusWriter{
			ResponseWriter: w,
			status:         200,
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		slog.Info("response",
			"request_id", requestID,
			"status", wrapped.status,
			"path", r.URL.Path,
			"method", r.Method,
			"duration_ms", duration.Milliseconds(),
			"slow", duration > time.Second*3,
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
