package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/luckyBambooBro/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpData struct {
	Content string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpSQL struct {
	ID uuid.UUID `json:"id"`
    Created_at time.Time `json:"createdAt"`
    Updated_at time.Time `json:"updatedAt"`
    Body string `json:"body"`
    User_id uuid.UUID `json:"user_id"`
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
	if chirp == nil {
		return
	}
	chirpParams := database.CreateChirpParams{
		Body: chirp.Content,
		UserID: chirp.UserID,
	}

	chirpCreate, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create chirp on database", err)
		return
	}
	chirpJSON := chirpSQL{
		ID: chirpCreate.ID,
		Created_at: chirpCreate.CreatedAt,
		Updated_at: chirpCreate.UpdatedAt,
		Body: chirpCreate.Body,
		User_id: chirp.UserID,
	}
	//if successfully created, return the values to user
	respondWithJSON(w, http.StatusCreated, chirpJSON)
}

func validateChirp(w http.ResponseWriter, chirp *chirpData) *chirpData {
	const maxChirpLength = 140

	//handle request depending on length
	if len(chirp.Content) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return nil
	}

	//filter and return valid chirp
	filterProfanities(chirp)
	return chirp
}

func filterProfanities(c *chirpData) {
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
	//update the contents of the chirpData
	c.Content = censoredText
}
