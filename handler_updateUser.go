package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
)

type UpdateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("invalid authecation: %s", err)
		w.WriteHeader(401)
		return
	}
	user_id, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		w.WriteHeader(401)
		return
	}
	decoder := json.NewDecoder(r.Body)
	myUser := &UpdateUser{}
	err = decoder.Decode(myUser)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		w.WriteHeader(400)
		return
	}
	data, err := auth.HashPassword(myUser.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}
	user, err := cfg.db.UpdateUserEmail(r.Context(), database.UpdateUserEmailParams{
		Email:          myUser.Email,
		HashedPassword: data,
		ID:             user_id,
	})
	if err != nil {
		log.Printf("Error updating email: %s", err)
		w.WriteHeader(500)
		return
	}
	responseUser := User{
		Email:     user.Email,
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	w.WriteHeader(200)

	err = encoder.Encode(responseUser)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
	}
}
