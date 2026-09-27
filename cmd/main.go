package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	api "github.com/nikbonk/go-http-server/internal/api"
	validateBody "github.com/nikbonk/go-http-server/internal/validateBody"
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

	cfg := api.NewApiConfig(db)

	appHandler := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(filepathRoot)),
	)

	mux := http.NewServeMux()
	mux.Handle("GET /app/", cfg.MiddlewareMetricsInc(api.MiddlewareLog(appHandler)))
	mux.Handle("GET /api/healthz", cfg.MiddlewareMetricsInc(api.MiddlewareLog(http.HandlerFunc(api.HealthHandler))))
	mux.Handle("GET /admin/metrics", api.MiddlewareLog(http.HandlerFunc(cfg.MetricHandler)))
	mux.Handle("POST /admin/reset", api.MiddlewareLog(http.HandlerFunc(cfg.ResetMetricsHandler)))
	mux.Handle("POST /api/validate_chirp", cfg.MiddlewareMetricsInc(api.MiddlewareLog(http.HandlerFunc(validateBody.ValidateBodyHandler))))

	srv := &http.Server{
		Addr:         ":" + srvPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, srvPort)
	log.Fatal(srv.ListenAndServe())

}
