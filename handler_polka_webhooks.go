package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/luckyBambooBro/chirpy/internal/auth"
)

const upgradedString = "user.upgraded"

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	//check api key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaAPI {
		respondWithError(w, http.StatusUnauthorized, "invalid or missing api key", err)
		return
	}

	//decode webhook req
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

	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), polkaWebhookRequest.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to update user to chirpy red", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)

}
