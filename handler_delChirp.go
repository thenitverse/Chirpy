package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error in authentication: %s", err)
		w.WriteHeader(401)
		return
	}
	tokenUser, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		w.WriteHeader(401)
		return
	}
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
	if onechirp.UserID != tokenUser {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     val,
		UserID: tokenUser,
	})
	if err != nil {
		log.Printf("Error deleting chirp: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
