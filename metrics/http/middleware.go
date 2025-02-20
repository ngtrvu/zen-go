package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode    int
	headerWritten bool
}

func NewCustomResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{w, http.StatusOK, false}
}

func (w *CustomResponseWriter) WriteHeader(code int) {
	if !w.headerWritten {
		w.StatusCode = code
		w.ResponseWriter.WriteHeader(code)
	}
}

func InboundMetricsMiddleware(observer HttpObserver) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var statusCode int
			defer func(time time.Time) {
				routePattern := chi.RouteContext(r.Context()).RoutePattern()

				// observe request
				observer.ObserveRequest(time, statusCode, r.Method, routePattern)
			}(time.Now())

			customRespWirter := NewCustomResponseWriter(w)
			next.ServeHTTP(customRespWirter, r)
			statusCode = customRespWirter.StatusCode
		}
		return http.HandlerFunc(fn)
	}
}
