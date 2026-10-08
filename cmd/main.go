package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	api "github.com/nikbonk/go-http-server/internal/api"
	"github.com/nikbonk/go-http-server/internal/database"
)

const (
	srvPort      = "8080"
	filepathRoot = "./html"
)

func loadApiConfig() (*api.ApiConfig, error) {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, errors.New("DB_URL not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	dbQueries := database.New(db)

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		return nil, errors.New("PLATFORM not set")
	} else if platform != "dev" && platform != "prod" {
		return nil, errors.New("PLATFORM must be dev or prod")
	}

	jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	if jwtSigningKey == "" {
		return nil, errors.New("JWT_SIGNING_KEY not set")
	}

	return api.NewApiConfig(dbQueries, platform, jwtSigningKey), nil
}

func main() {

	cfg, err := loadApiConfig()
	if err != nil {
		log.Fatal(err)
	}

	appHandler := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(filepathRoot)),
	)

	mux := http.NewServeMux()
	mux.Handle("GET /app/", cfg.MiddlewareMetricsInc(api.MiddlewareLog(appHandler)))

	mux.Handle("GET /admin/metrics", api.MiddlewareLog(http.HandlerFunc(cfg.MetricHandler)))
	mux.Handle("POST /admin/reset", api.MiddlewareLog(http.HandlerFunc(cfg.ResetHandler)))

	mux.Handle("GET /api/healthz", cfg.MiddlewareMetricsInc(api.MiddlewareLog(http.HandlerFunc(api.HealthHandler))))
	mux.Handle("POST /api/users", cfg.MiddlewareMetricsInc(api.MiddlewareLog(http.HandlerFunc(cfg.UserCreateHandler))))
	mux.Handle("POST /api/login", cfg.MiddlewareMetricsInc(api.MiddlewareLog(http.HandlerFunc(cfg.UserLoginHandler))))
	mux.Handle("POST /api/refresh", cfg.MiddlewareMetricsInc(api.MiddlewareLog(http.HandlerFunc(cfg.RefreshTokenHandler))))
	mux.Handle("POST /api/revoke", cfg.MiddlewareMetricsInc(api.MiddlewareLog(http.HandlerFunc(cfg.RevokeTokenHandler))))

	mux.Handle("POST /api/chirps", cfg.MiddlewareMetricsInc(api.MiddlewareLog(http.HandlerFunc(cfg.ChirpCreateHandler))))
	mux.Handle("GET /api/chirps", cfg.MiddlewareMetricsInc(api.MiddlewareLog(http.HandlerFunc(cfg.ChirpGetHandler))))
	mux.Handle("GET /api/chirps/{id}", cfg.MiddlewareMetricsInc(api.MiddlewareLog(http.HandlerFunc(cfg.ChirpGetByIdHandler))))

	srv := &http.Server{
		Addr:         ":" + srvPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, srvPort)
	log.Fatal(srv.ListenAndServe())

}
