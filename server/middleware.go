package server

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process the request
		next.ServeHTTP(lw, r)

		// Log request details
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lw.statusCode, duration)
	})
}

// loggingResponseWriter is a wrapper for http.ResponseWriter that captures the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}
