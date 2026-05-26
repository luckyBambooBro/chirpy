package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	//"github.com/luckyBambooBro/chirpy/internal/database"
)

type userData struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	//use requestData to decode data into go struct
	requestData, err := decodeUserInput(w, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request", err)
		return
	}
	
	//UP TO HERE: hash password before creating
	

	//create user
	createUser, err := cfg.db.CreateUser(r.Context(), requestData.Email)
	if err != nil {
		log.Printf("error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
		return
	}

	user := User{
		ID:        createUser.ID,
		CreatedAt: createUser.CreatedAt,
		UpdatedAt: createUser.UpdatedAt,
		Email:     createUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)
}

//come back to this after updating handlerUsersCreate
/*func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	//use requestData to decode data into go struct
	requestData, err := decodeUserInput(w, r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request", err)
		return
	}
}
*/

func decodeUserInput(w http.ResponseWriter, r *http.Request) (*userData, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	requestData := &userData{}
	if err := decoder.Decode(requestData); err != nil {
		log.Printf("error decoding user details: %v", err)
		return nil, err
	}
	return requestData, nil
}