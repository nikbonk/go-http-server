package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nikbonk/go-http-server/internal/auth"
	"github.com/nikbonk/go-http-server/internal/database"
)

func (cfg *ApiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerRefreshToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	foundToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "Unknown refresh token", http.StatusUnauthorized)
		return
	}

	if foundToken.ExpiresAt.Before(time.Now().UTC()) {
		http.Error(w, "Unknown refresh token", http.StatusUnauthorized)
		return
	}

	if foundToken.RevokedAt.Valid {
		http.Error(w, "Unknown refresh token", http.StatusUnauthorized)
		return
	}

	jwt, err := auth.MakeJWT(foundToken.UserID, cfg.jwtSigningKey)
	if err != nil {
		http.Error(w, "error while generating token", http.StatusInternalServerError)
		return
	}

	resp := refreshTokenResponse{
		JWT: jwt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "error while encoding response", http.StatusInternalServerError)
		return
	}
}

func (cfg *ApiConfig) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerRefreshToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	_, err = cfg.db.RevokeToken(r.Context(), database.RevokeTokenParams{
		RevokedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: now,
		Token:     refreshToken,
	})
	if err != nil {
		http.Error(w, "Unknown token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
