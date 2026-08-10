package main

import (
	"context"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	// Check for present query params
	authorID := r.URL.Query().Get("author_id")

	if authorID != "" {
		userID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}

		chirps, err := cfg.queries.GetChirpsByUser(context.Background(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve chirps", err)
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
		return
	}


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

	direction := r.URL.Query().Get("sort")
	if direction == "asc" {
		sort.Slice(responses, func(i, j int) bool { return responses[i].CreatedAt.Before(responses[j].CreatedAt) })
	} 
	if direction == "desc" {
		sort.Slice(responses, func(i, j int) bool { return responses[i].CreatedAt.After(responses[j].CreatedAt) })
	}

	respondWithJSON(w, http.StatusOK, responses)

}
