package main

import (
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Perform some processing before calling the next handler
		// log.Println("Received request:", r.URL)
		log.Printf("%s - %s", r.Method, r.URL)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Perform some processing after the next handler
		// log.Println("Completed request:", r.URL)
	})
}
