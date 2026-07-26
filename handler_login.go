package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) UserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	myUser := &createUserRequest{}
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	err = json.NewEncoder(w).Encode(responseUser)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
	}

}
