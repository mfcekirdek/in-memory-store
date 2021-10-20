package middleware

import (
	"log"
	"net/http"
)

type LoggerMiddleware struct {
	handler http.Handler
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := newLoggingResponseWriter(w)
	l.handler.ServeHTTP(lrw, r)
	statusCode := lrw.statusCode
	log.Printf("%s %s [%d] - %s", r.Method, r.URL.Path, statusCode, http.StatusText(statusCode))
}

func NewLoggerMiddleware(handlerToWrap http.Handler) *LoggerMiddleware {
	return &LoggerMiddleware{handler: handlerToWrap}
}
