package middleware

import "net/http"

// loggingResponseWriter wraps http.ResponseWriter interface and status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newLoggingResponseWriter creates a new loggingResponseWriter instance.
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

// WriteHeader overrides the WriteHeader method of the wrapped http.ResponseWriter instance and saves HTTP status code value.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
