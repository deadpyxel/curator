package main

import "net/http"

func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	type apiOk struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, apiOk{Status: http.StatusText(http.StatusOK)})
}

func errorEndpoint(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
