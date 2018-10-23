package chttp

import (
	"net/http"
)

// JSONResponseMiddleware is an HTTP-handling middleware that adds `Content-Type: application/json` to the response.
func JSONResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
