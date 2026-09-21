package main

import (
	"log"
	"net/http"
	"time"
)

const (
	srvPort      = "8080"
	filepathRoot = "."
)

func main() {

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:         ":" + srvPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	mux.Handle("/", http.FileServer(http.Dir(".")))

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, srvPort)
	log.Fatal(srv.ListenAndServe())

}
