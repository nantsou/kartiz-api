package main

import (
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (s *server) accessLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{w, http.StatusOK}
		next.ServeHTTP(lrw, r)
		log.Printf("\"%s %s %s\" %d %s", r.Method, r.RequestURI, r.Proto, lrw.statusCode, http.StatusText(lrw.statusCode))
	})
}

func (s *server) setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}