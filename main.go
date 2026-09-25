package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	srvPort      = "8080"
	filepathRoot = "./html"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	cfg := &apiConfig{}

	appHandler := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(filepathRoot)),
	)

	mux := http.NewServeMux()
	mux.Handle("GET /app/", cfg.middlewareMetricsInc(middlewareLog(appHandler)))
	mux.Handle("GET /api/healthz", cfg.middlewareMetricsInc(middlewareLog(http.HandlerFunc(healthHandler))))
	mux.Handle("GET /admin/metrics", middlewareLog(http.HandlerFunc(cfg.metricHandler)))
	mux.Handle("POST /admin/reset", middlewareLog(http.HandlerFunc(cfg.resetHandler)))
	mux.Handle("POST /api/validate_chirp", cfg.middlewareMetricsInc(middlewareLog(http.HandlerFunc(validateBodyHandler))))

	srv := &http.Server{
		Addr:         ":" + srvPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, srvPort)
	log.Fatal(srv.ListenAndServe())

}
