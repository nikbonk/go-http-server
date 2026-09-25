package http

import (
	"database/sql"
	"sync/atomic"
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
