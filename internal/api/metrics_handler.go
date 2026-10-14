package api

import (
	"html/template"
	"net/http"
)

// MetricHandler godoc
// @Summary      File server metrics
// @Description  Displays the file server hit counter as HTML.
// @Tags         admin
// @Produce      html
// @Success      200 {string} string "HTML metrics page"
// @Failure      500 {string} string "Internal server error"
// @Router       /admin/metrics [get]
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

// HealthHandler godoc
// @Summary      Health check
// @Description  Returns OK when the handler is reachable.
// @Tags         health
// @Produce      plain
// @Success      200 {string} string "OK"
// @Router       /api/healthz [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// imagine system check before just returning a 200 OK
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))
}
