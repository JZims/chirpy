package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JZims/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	inputParams := parameters{}
	err := decoder.Decode(&inputParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// If request provides ExpiresAt value under an hour, use it. Default to 1 hour

	expirationTime := time.Hour
	if inputParams.ExpiresInSeconds > 0 && inputParams.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(inputParams.ExpiresInSeconds) * time.Second
	}

	// Get User Info based on email
	user, err := cfg.queries.GetUserByEmail(ctx, inputParams.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	// Compare hashed passwords
	matchOk, err := auth.CheckPasswordHash(inputParams.Password, user.HashedPassword)
	if err != nil || !matchOk {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	// Use info to create JWT
	token, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create auth token", err)
		return
	}

	// Return user obj sans password
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})

}
