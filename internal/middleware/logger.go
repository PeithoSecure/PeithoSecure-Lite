package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware logs incoming HTTP requests in a simple, structured way.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Serve next handler
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("[%s] %s %s %s in %v", r.Method, r.RequestURI, r.RemoteAddr, r.UserAgent(), duration)
	})
}
