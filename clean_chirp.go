package main

import (
	"strings"
)

type cleanedResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func cleanChirp(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		currentWord := strings.ToLower(word)
		if currentWord == "kerfuffle" || currentWord == "sharbert" || currentWord == "fornax" {
			words[i] = "****"
		}

	}
	return strings.Join(words, " ")

}
