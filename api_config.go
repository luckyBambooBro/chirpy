package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	
}

func (cfg *apiConfig) apiHandler() {
	fmt.Printf("Hits: %v", cfg.fileserverHits)
}