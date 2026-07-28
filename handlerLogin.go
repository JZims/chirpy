package main

import (
	"encoding/json"
	"net/http"

	"github.com/JZims/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	inputParams := parameters{}
	err := decoder.Decode(&inputParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
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

	// Create JSON for response
	userObj := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	// Return user obj sans password
	respondWithJSON(w, http.StatusOK, userObj)

}
