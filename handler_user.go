package main

import (
	"encoding/json"

	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	myUser := &User{}
	err := decoder.Decode(myUser)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), myUser.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return

	}
	responseUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	w.WriteHeader(http.StatusCreated)

	err = encoder.Encode(responseUser)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
	}

}
