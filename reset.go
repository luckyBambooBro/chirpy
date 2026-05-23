package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "reset is only allowed in dev environment", nil)
		return
	}
	if err := cfg.db.DeleteAllUsers(r.Context()); err != nil {
		log.Printf("Failed to reset the database: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hits reset to 0 and database reset to initial state")
}
