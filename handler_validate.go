package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140

	type chirpData struct {
	Content string `json:"body"`
	}
	type chirpValid struct {
	ChirpValid bool `json:"valid"`
	}

	//decode request
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	chirp := &chirpData{}
	if err := decoder.Decode(chirp); err != nil {
		log.Printf("Error decoding request: %v", err)
		respondWithError(w, http.StatusInternalServerError, "unable to decode request")
		return
	}
	//handle request depending on length
	if len(chirp.Content) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	
	respondWithJSON(w, http.StatusOK, chirpValid{
		ChirpValid: true,
	})
}
