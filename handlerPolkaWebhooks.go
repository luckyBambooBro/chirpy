package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data struct {
		UserID string `json:"user_id"` 
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	upgradedString:= "user.upgraded"
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	polkaWebhookRequest := &polkaWebhookRequest{}
	if err := decoder.Decode(polkaWebhookRequest); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decodong polka webhook", err)
		return
	}

	if polkaWebhookRequest.Event != upgradedString {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	userUUID, err := uuid.Parse(polkaWebhookRequest.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to parse user ID", err)
		return
	}
	updatedUser, err := cfg.db.UpdateChirpyRed(r.Context(), userUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to update user to chirpy red", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, response{
		User{
			ID: updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email: updatedUser.Email,
			IsChirpyRed: updatedUser.IsChirpyRed,
		},
	})

}