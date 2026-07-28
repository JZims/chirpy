package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JZims/chirpy/internal/auth"
	"github.com/JZims/chirpy/internal/database"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Check for Password exist
	if len(params.Password) == 0 {
		err := errors.New("Password invalid or not provided")
		respondWithError(w, http.StatusBadRequest, "Password invalid or not provided", err)
		return
	}

	// Hash password before storing in db
	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash provided password", err)
		return
	}

	userObj := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPw,
	}

	newUser, err := cfg.queries.CreateUser(ctx, userObj)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create New User", err)
		return
	}

	mappedUser := User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}
	respondWithJSON(w, http.StatusCreated, mappedUser)

}
