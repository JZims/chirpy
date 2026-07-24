package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.queries.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch Chirps from Database.", err)
		return
	}

	responses := []returnChirpVals{}

	for _, c := range chirps {
		chirp := returnChirpVals{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
		responses = append(responses, chirp)
	}

	respondWithJSON(w, http.StatusOK, responses)

}
