package http

import (
	"database/sql"
	"net/http"
	"sync/atomic"
	"text/template"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	db             *sql.DB
}

func NewApiConfig(db *sql.DB) *ApiConfig {
	return &ApiConfig{
		db: db,
	}
}

func (cfg *ApiConfig) MetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	websiteHits := cfg.fileserverHits.Load()

	tmpl, err := template.ParseFiles("./html/metrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]int32{
		"websiteHits": websiteHits,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cfg *ApiConfig) ResetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
