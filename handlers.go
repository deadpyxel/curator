package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	respondWithJSON(w, http.StatusOK, dbUserToUser(dbUser))
}

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       params.Url,
		UserID:    dbUser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create new user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, dbFeedToFeed(feed))
}
