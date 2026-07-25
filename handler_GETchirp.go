package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	reqChirp, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
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
func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	val, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid UUID Parse: %s", err)
		w.WriteHeader(400)
		return
	}
	onechirp, err := cfg.db.GetChirp(r.Context(), val)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(404)
		return
	}
	if err != nil {
		log.Printf("Error getting chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	responseChirp := Chirp{
		ID:        onechirp.ID,
		CreatedAt: onechirp.CreatedAt,
		UpdatedAt: onechirp.UpdatedAt,
		Body:      onechirp.Body,
		UserID:    onechirp.UserID,
	}
	data, err := json.Marshal(responseChirp)
	if err != nil {
		log.Printf("Error marshalling data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)

}
