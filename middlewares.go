package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/deadpyxel/curator/internal/auth"
	"github.com/deadpyxel/curator/internal/database"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s - %s", r.Method, r.URL.Path, r.UserAgent())
		next.ServeHTTP(w, r)
	})
}

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) authMiddleware(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		handler(w, r, dbUser)
	}
}
