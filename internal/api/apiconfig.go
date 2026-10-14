package api

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/nikbonk/go-http-server/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSigningKey  string
	polkaKey       string
}

func NewApiConfig(db *database.Queries, platform string, jwtSigningKey string, polkaKey string) *ApiConfig {
	return &ApiConfig{
		db:            db,
		platform:      platform,
		jwtSigningKey: jwtSigningKey,
		polkaKey:      polkaKey,
	}
}

type chirpRequest struct {
	Body string `json:"body"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userCreateResponse struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	IsMembershipRed bool      `json:"is_chirpy_red"`
}

type userLoginResponse struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	JWT             string    `json:"token"`
	RefreshToken    string    `json:"refresh_token"`
	IsMembershipRed bool      `json:"is_chirpy_red"`
}

type refreshTokenResponse struct {
	JWT string `json:"token"`
}

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}
