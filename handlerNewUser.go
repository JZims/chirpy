package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	newUser, err := cfg.queries.CreateUser(ctx, params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create User", err)
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
