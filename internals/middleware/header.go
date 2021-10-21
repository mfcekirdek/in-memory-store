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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	l.handler.ServeHTTP(w, r)
}

func NewHeaderMiddleware(handlerToWrap http.Handler) *HeaderMiddleware {
	return &HeaderMiddleware{handlerToWrap}
}
