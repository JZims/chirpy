package main

import (
	"net/http"

	"github.com/JZims/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	err = cfg.queries.UpdateRefreshToken(r.Context(), reqToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update DB", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, response{})
}
