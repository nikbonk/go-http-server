package api

import (
	"sync/atomic"

	"github.com/nikbonk/go-http-server/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func NewApiConfig(db *database.Queries, platform string) *ApiConfig {
	return &ApiConfig{
		db:       db,
		platform: platform,
	}
}
