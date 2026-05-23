package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/luckyBambooBro/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpData struct {
	Content string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	chirp := &chirpData{}
	if err := decoder.Decode(chirp); err != nil {
		log.Printf("Error decoding request: %v", err)
		respondWithError(w, http.StatusInternalServerError, "unable to decode request\n", err)
		return
	}

	chirp = validateChirp(w, chirp)
	chirpParams := database.CreateChirpParams{
		Body: chirp.Content,
		UserID: chirp.UserID,
	}

	chirpSQL, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create chirp on database", err)
		return
	}

	//if successfully created, return the values to user
	respondWithJSON(w, http.StatusOK, chirpSQL)
}

func validateChirp(w http.ResponseWriter, chirp *chirpData) *chirpData {
	const maxChirpLength = 140

	//handle request depending on length
	if len(chirp.Content) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return nil
	}

	//filter and return valid chirp
	filteredChirp := filterProfanities(chirp)
	return filteredChirp
}

func filterProfanities(c *chirpData) *chirpData {
	//variables
	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(c.Content, " ")
	censor := "****"

	//filter
	for i, word := range words {
		wordLower := strings.ToLower(word)
		if _, ok := profanities[wordLower]; ok {
			words[i] = censor
		}
	}

	censoredText := strings.Join(words, " ")

	return &chirpData{
		Content: censoredText,
	}
}
