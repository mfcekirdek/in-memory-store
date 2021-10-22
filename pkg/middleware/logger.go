package middleware

import (
	"log"
	"net/http"
)

// LoggerMiddleware wraps http.Handler interface.
type LoggerMiddleware struct {
	handler http.Handler
}

// ServeHTTP logs information according to HTTP requests after they are executed.
func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := newLoggingResponseWriter(w)
	l.handler.ServeHTTP(lrw, r)
	statusCode := lrw.statusCode
	log.Printf("%s %s [%d] - %s", r.Method, r.URL.Path, statusCode, http.StatusText(statusCode))
}

// NewLoggerMiddleware creates a new LoggerMiddleware instance.
func NewLoggerMiddleware(handlerToWrap http.Handler) *LoggerMiddleware {
	return &LoggerMiddleware{handler: handlerToWrap}
}
