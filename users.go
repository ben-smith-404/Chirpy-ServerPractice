package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ben-smith-404/Chirpy-ServerPractice/internal/auth"
	"github.com/ben-smith-404/Chirpy-ServerPractice/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type AuthUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handleAddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	newUser := AuthUser{}
	err := decoder.Decode(&newUser)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, 500, "something went wrong")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          newUser.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, 500, "something went wrong")
		return
	}
	ReturnUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, 201, ReturnUser)
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	authUser := AuthUser{}
	err := decoder.Decode(&authUser)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "something went wrong")
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), authUser.Email)
	if err != nil {
		log.Printf("Error retrieving user from database: %s", err)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	passwordMatch, err := auth.CheckPasswordHash(authUser.Password, dbUser.HashedPassword)
	if err != nil || !passwordMatch {
		log.Printf("Error matching password to hashed password: %s", err)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	ReturnUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJSON(w, 200, ReturnUser)
}
