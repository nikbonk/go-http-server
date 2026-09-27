package api

import "net/http"

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Deleting users on production is not allowed"))
		return
	}
	cfg.fileserverHits.Store(0)
	if err := cfg.db.ResetUsers(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset users"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users reset successfully"))
}
