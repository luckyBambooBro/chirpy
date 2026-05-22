package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	//"github.com/luckyBambooBro/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type userData struct{
		Email string `json:"email"`
	}
	
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	requestData := &userData{}
	if err := decoder.Decode(requestData); err != nil {

		log.Printf("error decoding request to create user: %v", err)
		returnErrorMsg(w, http.StatusInternalServerError, "internal server error")
		return
	}
	//create user
	createUser, err := cfg.db.CreateUser(r.Context(), requestData.Email)
	if err != nil {
		log.Printf("error creating user: %v", err)
		returnErrorMsg(w, http.StatusInternalServerError,  "internal server error")
		return
	}

	user := User{
		ID: createUser.ID,
		CreatedAt: createUser.CreatedAt,
		UpdatedAt: createUser.UpdatedAt,
		Email: createUser.Email,
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("error encoding user data: %v", err)
		return
	}

	//return data in json object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(data); err != nil {
		log.Println(err)
	}

}


