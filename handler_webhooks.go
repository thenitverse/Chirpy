package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type webhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) webhookReq(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	myWebhook := &webhookRequest{}
	err := decoder.Decode(myWebhook)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}
	if myWebhook.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(myWebhook.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = cfg.db.UpdateChirpy(r.Context(), userID)
	if err != nil {
		log.Printf("Error upgrading chirpy: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
