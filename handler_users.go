package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/luckyBambooBro/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type userData struct{
		userEmail string
	}
	
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	requestData := &userData{}
	if err := decoder.Decode(requestData); err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
		log.Printf("error decoding request to create user: %v", err)
		return
	}
	//come back to the following after decoding the resp
	cfg.db.CreateUser(r.Context(), requestData.userEmail)
}
