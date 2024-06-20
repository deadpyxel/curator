package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/deadpyxel/curator/internal/auth"
	"github.com/deadpyxel/curator/internal/database"
	"github.com/google/uuid"
)

func handlerLiveness(w http.ResponseWriter, r *http.Request) {
	type apiOk struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, apiOk{Status: http.StatusText(http.StatusOK)})
}

func handlerErrorTest(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create new user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, dbUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth Error: %v", err))
		return
	}

	dbUser, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failure to fetch user information: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbUserToUser(dbUser))
}
