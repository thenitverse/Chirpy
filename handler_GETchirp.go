package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	reqChirp, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirp: %s", err)
		w.WriteHeader(500)
		return
	}
	Chirps := []Chirp{}
	for _, item := range reqChirp {
		respChirp := Chirp{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Body:      item.Body,
			UserID:    item.UserID,
		}
		Chirps = append(Chirps, respChirp)

	}
	data, err := json.Marshal(Chirps)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)

}
