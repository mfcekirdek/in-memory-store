package middleware

import (
	"log"
	"net/http"
)

type LoggerMiddleware struct {
	handler http.Handler
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s", r.Method, r.URL.Path)
}

func NewLoggerMiddleware(handlerToWrap http.Handler) *LoggerMiddleware {
	return &LoggerMiddleware{handlerToWrap}
}
