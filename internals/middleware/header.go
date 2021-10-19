package middleware

import (
	"net/http"
)

type HeaderMiddleware struct {
	handler http.Handler
}

func (l *HeaderMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	l.handler.ServeHTTP(w, r)
}

func NewHeaderMiddleware(handlerToWrap http.Handler) *HeaderMiddleware {
	return &HeaderMiddleware{handlerToWrap}
}
