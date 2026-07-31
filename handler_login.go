package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) UserLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type loginResponse struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)
	myUser := &loginRequest{}
	err := decoder.Decode(myUser)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), myUser.Email)
	if err != nil {
		log.Printf("Unauthorized: %s", err)
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}
	matches, err := auth.CheckPasswordHash(
		myUser.Password,
		user.HashedPassword,
	)
	if !matches || err != nil {
		log.Printf("401 Unauthorized: %s", err)
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return

	}
	expiresIn := time.Hour

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		log.Printf("Error creating token: %s", err)
		w.WriteHeader(500)
		return
	}
	refreshToken := auth.MakeRefreshToken()
	now := time.Now()
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		CreatedAt: now,
		UpdatedAt: now,
		RevokedAt: sql.NullTime{},
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		log.Printf("Data constraint failed: %s", err)
		w.WriteHeader(500)
		return
	}
	responseUser := loginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseUser)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
	}

}
