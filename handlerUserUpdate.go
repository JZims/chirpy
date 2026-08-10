package main

import (
	"encoding/json"
	"net/http"

	"github.com/JZims/chirpy/internal/auth"
	"github.com/JZims/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {

	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	inputParams := reqParams{}
	err := decoder.Decode(&inputParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	hashedPassword, err := auth.HashPassword(inputParams.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't encode password", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	u, err := cfg.queries.GetUserById(r.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	updatedUser, err := cfg.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          inputParams.Email,
		HashedPassword: hashedPassword,
		ID:             u.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          updatedUser.ID,
		Email:       updatedUser.Email,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})

}
