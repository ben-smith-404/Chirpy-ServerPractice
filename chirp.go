package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handleChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	type filteredString struct {
		CleanedBody string `json:"cleaned_body"`
	}
	responseBody := filteredString{
		CleanedBody: maskProfanity(params.Body),
	}
	respondWithJSON(w, 200, responseBody)
}

func maskProfanity(chirp string) string {
	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(chirp, " ")
	for key, word := range words {
		if slices.Contains(forbiddenWords, strings.ToLower(word)) {
			words[key] = "****"
		}
	}
	return strings.Join(words, " ")
}
