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
	IsChirpyRed    bool	`json:"is_chirpy_red"`
}

type response struct {
	User
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

//constants
const (
	jwtExpiry = 1 * time.Hour
	refreshTokenExpiry = 60 * 24 * time.Hour
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	//use userData to decode data into go struct
	userData, err := decodeUserInput(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding user request", err)
		return
	}
	
	//hash the provided password before creating
	hashedPassword, err := auth.HashPassword(userData.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}
	//create user
	userParams := database.CreateUserParams{
		Email: userData.Email,
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
		IsChirpyRed: false,
	}

	respondWithJSON(w, http.StatusCreated, response{User: user})
}


func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	//use userData to decode data into go struct
	userData, err := decodeUserInput(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding user request", err)
		return
	}

	//obtain the user
	user, err := cfg.db.GetUserByEmail(r.Context(), userData.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	
	//compare passwords
	correctPassword, err := auth.CheckPasswordHash(userData.Password, user.HashedPassword)
	if !correctPassword {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	//create jwt
	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, jwtExpiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create access token", err)
		return
	}

	//create refresh token
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create refresh token", nil)
		return
	}
	//create token on database here
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token: refreshTokenString,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refreshToken on database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, 
		response{
			User: User{
				ID: user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email: user.Email,
				IsChirpyRed: user.IsChirpyRed,
				},
				Token: jwt,
				RefreshToken: refreshTokenString,
		},
	)

}


func (cfg *apiConfig) handlerUserRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized access", err)
		return
	}
	if len(refreshToken) != 64 { //apparently len is 64 instead of 32 due to hex encoding (says gemini)
		respondWithError(w, http.StatusUnauthorized, "unauthorized access", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid access", err)
		return
	}

	
	//if refresh token valid provide access token 
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, jwtExpiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
	
}

func (cfg *apiConfig) handlerUserRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized access", err)
		return
	}

	if _, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error revoking refresh token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerUpdateEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	userData, err := decodeUserInput(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding user request", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error retrieving bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorised access", err)
		return
	}

	//hash password into database
	hashedPassword, err := auth.HashPassword(userData.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	updatedUser, err := cfg.db.UpdateEmailAndPassword(r.Context(), database.UpdateEmailAndPasswordParams{
		HashedPassword: hashedPassword,
		Email: userData.Email,
		ID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to update user details to database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email: updatedUser.Email,
		},
		Token: accessToken, //gemini: usually you would make a new JWT and send it back (more secure) but lesson doesnt ask for this
	})

	
}

//==========HELPER FUNCTIONS =================
func decodeUserInput(r *http.Request) (*userData, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	userData := &userData{}
	if err := decoder.Decode(userData); err != nil {
		log.Printf("error decoding user request: %v", err)
		return nil, err
	}
	return userData, nil
}