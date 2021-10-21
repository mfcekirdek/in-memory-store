package middleware

import (
	"net/http"
)

type HeaderMiddleware struct {
	handler http.Handler
}

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

func NewHeaderMiddleware(handlerToWrap http.Handler) *HeaderMiddleware {
	return &HeaderMiddleware{handlerToWrap}
}
