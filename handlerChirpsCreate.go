package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JZims/chirpy/internal/auth"
	"github.com/JZims/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Validatae User has a JWT
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials, bruv", err)
		return
	}

	log.Printf("%v", token)

	userRequesting, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials, innit", err)
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
		UserID:    userRequesting,
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
