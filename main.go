package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const srvPort = "8080"

func main() {

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:         ":" + srvPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Printf("Server listening on %v\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}
