package main

import (
	"log"
	"net/http"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s - %s", r.Method, r.URL.Path, r.UserAgent())
		next.ServeHTTP(w, r)
	})
}
