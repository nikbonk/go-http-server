package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nikbonk/go-http-server/internal/auth"
	"github.com/nikbonk/go-http-server/internal/database"
)

type userCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type user struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *ApiConfig) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req userCreateRequest
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "error while decoding request", http.StatusInternalServerError)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "error while hashing password", http.StatusInternalServerError)
		return
	}

	createdUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		http.Error(w, "error while creating user", http.StatusInternalServerError)
		return
	}

	// Do NOT return the hashed password in the respone (naughty naughty!)
	resp := user{
		ID:        createdUser.ID,
		Email:     req.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "error while encoding response", http.StatusInternalServerError)
		return
	}

}

func (cfg *ApiConfig) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req userCreateRequest
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "error while decoding request", http.StatusInternalServerError)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	requestedUser, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "error while fetching user", http.StatusInternalServerError)
		return
	}

	passwordMatches, err := auth.CheckPasswordHash(req.Password, requestedUser.HashedPassword)
	if err != nil {
		http.Error(w, "Error while checking password", http.StatusInternalServerError)
		return
	}

	if !passwordMatches {
		http.Error(w, "password does not match", http.StatusUnauthorized)
		return
	}

	// Do NOT return the hashed password in the respone (naughty naughty!)
	resp := user{
		ID:        requestedUser.ID,
		Email:     req.Email,
		CreatedAt: requestedUser.CreatedAt,
		UpdatedAt: requestedUser.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "error while encoding response", http.StatusInternalServerError)
		return
	}
}
