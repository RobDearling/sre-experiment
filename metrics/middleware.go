package metrics

import (
	"net/http"
	"strconv"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ConcurrentRequests.Inc()
		defer ConcurrentRequests.Dec()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrapped.statusCode)
		RequestDuration.WithLabelValues(
			r.URL.Path,
			r.Method,
			status,
		).Observe(duration)

		RequestsTotal.WithLabelValues(
			r.URL.Path,
			r.Method,
			status,
		).Inc()

		if wrapped.statusCode >= 400 {
			ErrorsTotal.WithLabelValues(
				r.URL.Path,
				r.Method,
				status,
			).Inc()
		}
	})
}
