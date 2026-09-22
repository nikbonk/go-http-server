package main

import (
	"log"
	"net/http"
)

func appHandler(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/app", http.FileServer(http.Dir("."))).ServeHTTP(w, r)

	// not logging the response here, would need to wrap the response writer
	log.Printf("Request: %v", r.URL)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// imagine system check before just returning a 200 OK
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))

	log.Printf("Request: %v Response: %v", r.URL, http.StatusText(http.StatusOK))
}
