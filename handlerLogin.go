package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JZims/chirpy/internal/auth"
	"github.com/JZims/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	now := time.Now().UTC()

	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	inputParams := parameters{}
	err := decoder.Decode(&inputParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	expirationTime := time.Hour

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

	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.queries.CreateRefresh(r.Context(), database.CreateRefreshParams{
		Token:     refreshToken,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create refresh token", err)
		return
	}

	// Return user obj sans password, JWT and Refresh
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})

}
