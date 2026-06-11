package main

import (
	"encoding/json"
	"net/http"
)

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data struct {
		UserID string `json:"user_id"` 
	} `json:"data"`
}

func handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	upgradedString:= "user.upgraded"
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	polkaWebhookRequest := &polkaWebhookRequest{}
	if err := decoder.Decode(polkaWebhookRequest); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decodong polka webhook", err)
	}

	if polkaWebhookRequest.Event != upgradedString {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}


}