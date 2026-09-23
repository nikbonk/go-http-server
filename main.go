package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	srvPort      = "8080"
	filepathRoot = "."
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	cfg := &apiConfig{}

	appHandler := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(".")),
	)

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(middlewareLog(appHandler)))
	mux.Handle("/healthz", cfg.middlewareMetricsInc(middlewareLog(http.HandlerFunc(healthHandler))))
	mux.Handle("/metrics", middlewareLog(http.HandlerFunc(cfg.metricHandler)))
	mux.Handle("/reset", middlewareLog(http.HandlerFunc(cfg.resetHandler)))

	srv := &http.Server{
		Addr:         ":" + srvPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, srvPort)
	log.Fatal(srv.ListenAndServe())

}
