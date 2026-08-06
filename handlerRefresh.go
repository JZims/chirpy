package main

import (
	"net/http"
	"time"

	"github.com/JZims/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	userId, err := cfg.queries.GetUserFromRefresh(r.Context(), reqToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	newToken, err := auth.MakeJWT(userId, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create auth token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})

}
