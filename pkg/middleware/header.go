// Package middleware wraps handlers to do some pre- and/or post-processing of the request.
package middleware

import (
	"net/http"
)

// HeaderMiddleware wraps http.Handler interface.
type HeaderMiddleware struct {
	handler http.Handler
}

// ServeHTTP adds common headers to all requests before they are executed.
func (l *HeaderMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}
	l.handler.ServeHTTP(w, r)
}

// NewHeaderMiddleware creates a new HeaderMiddleware instance.
func NewHeaderMiddleware(handlerToWrap http.Handler) *HeaderMiddleware {
	return &HeaderMiddleware{handlerToWrap}
}
