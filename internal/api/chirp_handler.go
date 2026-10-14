package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/nikbonk/go-http-server/internal/auth"
	"github.com/nikbonk/go-http-server/internal/database"
	validate "github.com/nikbonk/go-http-server/internal/validateBody"
)

// ChirpCreateHandler godoc
// @Summary      Create chirp
// @Tags         chirps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body chirpRequest true "Chirp body"
// @Success      201 {object} chirpResponse
// @Failure      400 {string} string "Invalid JSON or chirp body"
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {string} string "Internal server error"
// @Router       /api/chirps [post]
func (c *ApiConfig) ChirpCreateHandler(w http.ResponseWriter, r *http.Request) {
	request := chirpRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	userID, err := auth.ValidateJWT(token, c.jwtSigningKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	valid, body, err := validate.ValidaChirp(request.Body)
	if !valid {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	chirp, err := c.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Body:      body,
		UserID:    userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// ChirpGetHandler godoc
// @Summary      List chirps
// @Description  Optionally filter chirps by author UUID.
// @Tags         chirps
// @Produce      json
// @Param        author_id query string false "Author UUID"
// @Success      200 {array} chirpResponse
// @Failure      400 {string} string "Invalid author_id"
// @Failure      500 {string} string "Internal server error"
// @Router       /api/chirps [get]
func (c *ApiConfig) ChirpGetHandler(w http.ResponseWriter, r *http.Request) {
	var (
		chirps []database.Chirp
		err    error
	)

	query := r.URL.Query().Get("author_id")

	if query == "" {
		chirps, err = c.db.GetChirps(r.Context())
	} else {
		authorID, parseErr := uuid.Parse(query)
		if parseErr != nil {
			http.Error(w, "Invalid author_id", http.StatusBadRequest)
			return
		}

		chirps, err = c.db.GetChirpsByAuthor(r.Context(), authorID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]chirpResponse, len(chirps))
	for i, chirp := range chirps {
		responses[i] = chirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ChirpGetByIdHandler godoc
// @Summary      Get chirp by ID
// @Tags         chirps
// @Produce      json
// @Param        id path string true "Chirp UUID"
// @Success      200 {object} chirpResponse
// @Failure      400 {string} string "Invalid chirp ID"
// @Failure      404 {string} string "Chirp not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /api/chirps/{id} [get]
func (c *ApiConfig) ChirpGetByIdHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid chirp ID", http.StatusBadRequest)
		return
	}
	chirp, err := c.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "chirp not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ChirpDeleteByIdHandler godoc
// @Summary      Delete chirp
// @Description  Only the author of a chirp may delete it.
// @Tags         chirps
// @Security     BearerAuth
// @Param        id path string true "Chirp UUID"
// @Success      204 "Chirp deleted"
// @Failure      400 {string} string "Invalid chirp ID"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "Not allowed to delete this chirp"
// @Failure      404 {string} string "Chirp not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /api/chirps/{id} [delete]
func (c *ApiConfig) ChirpDeleteByIdHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid chirp ID", http.StatusBadRequest)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	requestorID, err := auth.ValidateJWT(token, c.jwtSigningKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chirp, err := c.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "chirp not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if requestorID != chirp.UserID {
		http.Error(w, "Not allowed to delete this chirp", http.StatusForbidden)
		return
	}

	err = c.db.DeleteChirpByID(r.Context(), chirp.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
