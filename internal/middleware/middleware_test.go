package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	handler := RecoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))

	handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code \"ok\", got %d", w.Code)
	}

	if !bytes.Contains(buf.Bytes(), []byte("request_id")) {
		t.Error("expected request_id in the log")
	}

	if !bytes.Contains(buf.Bytes(), []byte("status")) {
		t.Error("expected status in the log")
	}

	if !bytes.Contains(buf.Bytes(), []byte("path")) {
		t.Error("expected path in the log")
	}

	if !bytes.Contains(buf.Bytes(), []byte("method")) {
		t.Error("expected method in the log")
	}

	if !bytes.Contains(buf.Bytes(), []byte("duration_ms")) {
		t.Error("expected duration_ms in the log")
	}

	if !bytes.Contains(buf.Bytes(), []byte("slow")) {
		t.Error("expected slow in the log")
	}
}
