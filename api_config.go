package main

import (
	"sync/atomic"

	"github.com/luckyBambooBro/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}
