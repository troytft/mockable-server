package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggerHandler(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/path?query", strings.NewReader("OK"))
	w := httptest.NewRecorder()

	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)

	noopHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	handler := NewLoggerHandler(logger, noopHandler)
	handler.ServeHTTP(w, r)

	got := w.Result()
	if got == nil {
		t.Fatalf("unexpected Result")
	}

	line := strings.TrimSpace(buffer.String())
	if line != `{"Method":"POST","URI":"/path?query","Headers":{},"Body":"OK"}` {
		t.Fatalf("unexpected log: `%s`", line)
	}
}