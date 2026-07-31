package main

import (
	"chirpy/internal/auth"
	"log"
	"net/http"
)

func (cfg *apiConfig) userRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error revoking token: %s", err)
		w.WriteHeader(401)
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf("Error finding the user: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
