package http

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}
