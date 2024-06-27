package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/deadpyxel/curator/internal/database"
	"github.com/google/uuid"
)

// handlerLiveness responds with a JSON message containing the status of the API
func handlerLiveness(w http.ResponseWriter, r *http.Request) {
	type apiOk struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, apiOk{Status: http.StatusText(http.StatusOK)})
}

// handlerErrorTest is a test endpoint to verify that error messages can be shown to the client.
func handlerErrorTest(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	// Decode JSON contents for processing
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// Create new user on database
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create new feed: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    dbUser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error occurred when creating feed follow: %v", err))
		return
	}

	rFeed := dbFeedToFeed(feed)
	rFeedFollow := dbFeedFollowToFeedFollow(feedFollow)
	response := struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}{
		Feed:       rFeed,
		FeedFollow: rFeedFollow,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "No feeds found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not retrieve feeds: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbFeedsToFeeds(feeds))
}

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    dbUser.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not follow feed: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, dbFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollowForUser(r.Context(), dbUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "The user is not following any feeds")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve feeds followed by the user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	feedFollowIDStr := r.PathValue("feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing feed follow ID: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: dbUser.ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Specified feed not found or not followed by user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete feed follow: %v", err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}

func (apiCfg *apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	queryLimit := 10 // set default query limit to 10 posts
	// If limit was specified as a query param, use that limit instead
	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Unable to parse limit query param")
			return
		}
		queryLimit = parsedLimit
	}

	posts, err := apiCfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: dbUser.ID,
		Limit:  int32(queryLimit),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "No posts found for this user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve posts: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbPostsToPosts(posts))
}
