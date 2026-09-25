package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	nbhttp "github.com/nikbonk/go-http-server/internal/http"
)

const (
	srvPort      = "8080"
	filepathRoot = "./html"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	cfg := nbhttp.NewApiConfig(db)

	appHandler := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(filepathRoot)),
	)

	mux := http.NewServeMux()
	mux.Handle("GET /app/", cfg.MiddlewareMetricsInc(nbhttp.MiddlewareLog(appHandler)))
	mux.Handle("GET /api/healthz", cfg.MiddlewareMetricsInc(nbhttp.MiddlewareLog(http.HandlerFunc(nbhttp.HealthHandler))))
	mux.Handle("GET /admin/metrics", nbhttp.MiddlewareLog(http.HandlerFunc(cfg.MetricHandler)))
	mux.Handle("POST /admin/reset", nbhttp.MiddlewareLog(http.HandlerFunc(cfg.ResetHandler)))
	mux.Handle("POST /api/validate_chirp", cfg.MiddlewareMetricsInc(nbhttp.MiddlewareLog(http.HandlerFunc(nbhttp.ValidateBodyHandler))))

	srv := &http.Server{
		Addr:         ":" + srvPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, srvPort)
	log.Fatal(srv.ListenAndServe())

}
