package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type response struct {
	AccessToken string `json:"token"`
}

func (cfg *apiConfig) userRefresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error refreshing token: %s", err)
		w.WriteHeader(401)
		return
	}
	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("Error finding user: %s", err)
		w.WriteHeader(401)
		return
	}
	accessToken, err := auth.MakeJWT(userID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error accessing token: %s", err)
		w.WriteHeader(500)
		return
	}
	responseToken := response{
		AccessToken: accessToken,
	}
	data, err := json.Marshal(responseToken)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)

}
