package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JZims/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string    `json:"body"`
		User uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// Send body to profanity filter
	filtered := profanityFilter(params.Body)

	// Map the request params to returnVals
	newChirp := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      filtered,
		UserID:    params.User,
	}

	// Create Chirp in database
	postedChirp, err := cfg.queries.CreateChirp(r.Context(), newChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create record in database", err)
		return
	}

	returnChirp := returnChirpVals{
		ID:        postedChirp.ID,
		CreatedAt: postedChirp.CreatedAt,
		UpdatedAt: postedChirp.UpdatedAt,
		Body:      postedChirp.Body,
		UserID:    postedChirp.UserID,
	}

	// Respond with JSON payload and 201
	respondWithJSON(w, http.StatusCreated, returnChirp)

}
