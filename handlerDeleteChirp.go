package main

import (
	"net/http"

	"github.com/JZims/chirpy/internal/auth"
	"github.com/JZims/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	type response struct {
		ChirpID uuid.UUID `json:"chirp_id"`
	}

	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	reqUserId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid Auth Token", err)
		return
	}

	chirp, err := cfg.queries.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No chirp found for given ID", err)
		return
	}

	if chirp.UserID != reqUserId {
		respondWithError(w, http.StatusForbidden, "Not Authorized to access chirp", err)
		return
	}

	err = cfg.queries.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpId,
		UserID: reqUserId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
