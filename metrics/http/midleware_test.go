package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert"
	metrics_http "github.com/ngtrvu/zen-go/metrics/http"
)

type mockResponseWriter struct {
	http.ResponseWriter
	writeHeaderCalls int
	lastStatusCode   int
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.writeHeaderCalls++
	m.lastStatusCode = statusCode
}

func TestCustomResponseWriter_WriteHeader(t *testing.T) {
	t.Run("single WriteHeader call", func(t *testing.T) {
		mock := &mockResponseWriter{}
		writer := metrics_http.NewCustomResponseWriter(mock)

		writer.WriteHeader(http.StatusBadRequest)

		assert.Equal(t, http.StatusBadRequest, writer.StatusCode)
		assert.Equal(t, 1, mock.writeHeaderCalls)
		assert.Equal(t, http.StatusBadRequest, mock.lastStatusCode)
	})

	t.Run("multiple WriteHeader calls", func(t *testing.T) {
		mock := &mockResponseWriter{}
		writer := metrics_http.NewCustomResponseWriter(mock)

		// First call
		writer.WriteHeader(http.StatusBadRequest)

		// Second call - should trigger the superfluous warning
		writer.WriteHeader(http.StatusNotFound)

		// Verify multiple calls occurred (reproducing the bug)
		assert.Equal(t, 2, mock.writeHeaderCalls)

		// Verify last status was written (showing the issue)
		assert.Equal(t, http.StatusNotFound, mock.lastStatusCode)
	})

	t.Run("integration with http.Handler", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)

			// Simulate middleware or error handler writing status again
			w.WriteHeader(http.StatusInternalServerError)
		})

		rec := httptest.NewRecorder()
		writer := metrics_http.NewCustomResponseWriter(rec)
		handler.ServeHTTP(writer, httptest.NewRequest("GET", "/", nil))

		// Shouldn't be allowed to override status code
		assert.Equal(t, http.StatusInternalServerError, writer.StatusCode)

		// Response should show both writes occurred
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
