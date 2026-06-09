package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/luckyBambooBro/chirpy/internal/auth"
	"github.com/luckyBambooBro/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpData struct {
	Content string `json:"body"`
}

type chirpSQL struct {
	ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    Body string `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	
	//decode request
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	chirp := &chirpData{}
	if err := decoder.Decode(chirp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode request\n", err)
		return
	}

	//authenticate user via JWT
	userJWT, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to authenticate user", err)
		return
	}

	userID, err := auth.ValidateJWT(userJWT, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to authenticate user", err)
		return
	}

	//create chirp
	chirp = validateChirp(w, chirp)
	if chirp == nil {
		return
	}
	chirpParams := database.CreateChirpParams{
		Body: chirp.Content,
		UserID: userID,
	}

	chirpCreate, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create chirp on database", err)
		return
	}

	//respond to request
	chirpJSON := chirpSQL{
		ID: chirpCreate.ID,
		CreatedAt: chirpCreate.CreatedAt,
		UpdatedAt: chirpCreate.UpdatedAt,
		Body: chirpCreate.Body,
		UserID: userID,
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

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpsList, err := cfg.db.GetChirpsAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to obtain chirps list", err)
		return
	}

	chirpsJSONList := []chirpSQL{}
	for _, chirp := range chirpsList {
		chirpsJSONList = append(chirpsJSONList, chirpSQL{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsJSONList)
}

func (cfg *apiConfig) handlerChirpsGetOne(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirpID provided", nil)
		return
	}
	

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to retrieve chirp from database", err)
		return
	}

	chirpJSON := chirpSQL{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirpJSON)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	
}