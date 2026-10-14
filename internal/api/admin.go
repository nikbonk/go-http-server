package api

import "net/http"

// ResetHandler godoc
// @Summary      Reset users
// @Description  Deletes users and resets the file server hit counter. Only allowed in dev.
// @Tags         admin
// @Produce      plain
// @Success      200 {string} string "Users reset successfully"
// @Failure      403 {string} string "Deleting users on production is not allowed"
// @Failure      500 {string} string "Failed to reset users"
// @Router       /admin/reset [post]
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
