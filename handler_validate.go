package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type chirpData struct {
	Content string `json:"body"`
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140


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

func filterProfanities (c chirpData) chirpData {
	//variables
	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	words := strings.Split(c.Content, " ")
	censor := "****"

	for i, word := range words {
		wordLower := strings.ToLower(word)
		if _, ok := profanities[wordLower]; ok {
			words[i] = censor
		}
	}

	censoredText := strings.Join(words, " ")

	return chirpData{
		Content: censoredText,
	}

}
