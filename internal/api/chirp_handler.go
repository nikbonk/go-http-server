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
