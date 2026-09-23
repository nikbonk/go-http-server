package main

import (
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// imagine system check before just returning a 200 OK
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))

}
