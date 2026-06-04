package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/luckyBambooBro/chirpy/internal/auth"
	"github.com/luckyBambooBro/chirpy/internal/database"
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
	Password  string	`json:"-"`
}

type response struct {
	User
	AccessToken string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	//use requestData to decode data into go struct
	requestData, err := decodeUserInput(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding user details", err)
		return
	}
	
	//hash the provided password before creating
	hashedPassword, err := auth.HashPassword(requestData.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}
	//create user
	userParams := database.CreateUserParams{
		Email: requestData.Email,
		HashedPassword: hashedPassword,
	}
	createUser, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
		return
	}

	user := User{
		ID:        createUser.ID,
		CreatedAt: createUser.CreatedAt,
		UpdatedAt: createUser.UpdatedAt,
		Email:     createUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, response{User: user})
}


func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	//use requestData to decode data into go struct
	requestData, err := decodeUserInput(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request", err)
		return
	}

	//obtain the user
	user, err := cfg.db.GetUserByEmail(r.Context(), requestData.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	
	//compare passwords
	correctPassword, err := auth.CheckPasswordHash(requestData.Password, user.HashedPassword)
	if !correctPassword {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	//create jwt
	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, 1 * time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create access token", err)
	}

	//CREATE REFRESH TOKEN HERE 
	refreshToken := auth.MakeRefreshToken()
	if refreshToken == "" {
		respondWithError(w, http.StatusInternalServerError, "unable to create refresh token", nil)
	}
	//create token on database here

	respondWithJSON(w, http.StatusOK, 
		response{
			User: User{
				ID: user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email: user.Email,
				},
				AccessToken: jwt,
		},
	)

}


func decodeUserInput(r *http.Request) (*userData, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	requestData := &userData{}
	if err := decoder.Decode(requestData); err != nil {
		log.Printf("error decoding user details: %v", err)
		return nil, err
	}
	return requestData, nil
}