package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/JZims/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateChirpyRed(w http.ResponseWriter, r *http.Request) {

	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data
	}

	decoder := json.NewDecoder(r.Body)
	inputParams := parameters{}
	err := decoder.Decode(&inputParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to marshal input parameters", err)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No API Key Found", err)
		return
	}

	if apiKey != os.Getenv("POLKA_KEY") {
		fmt.Printf("%v", apiKey)
		respondWithError(w, http.StatusUnauthorized, "Invalid API Key", err)
		return
	}

	if inputParams.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Event type not supported", err)
		return
	}

	err = cfg.queries.UpdateChirpyRed(r.Context(), inputParams.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to update user field", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
