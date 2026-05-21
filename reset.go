package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	if err := cfg.db.DeleteAllUsers(r.Context()); err != nil {
		log.Printf("error deleting users: %v", err)
		returnErrorMsg(w, http.StatusInternalServerError)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d\n", cfg.fileserverHits.Load())
}
