package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func MakeRefreshToken() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("Error making bytes: %s", err)
		return ""
	}

	encodedStr := hex.EncodeToString(key)
	return encodedStr

}
