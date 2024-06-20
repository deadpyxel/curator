package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	type apiOk struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, apiOk{Status: http.StatusText(http.StatusOK)})
}

func errorEndpoint(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", readinessEndpoint)
	mux.HandleFunc("GET /v1/err", errorEndpoint)
	logMux := logMiddleware(mux)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: logMux,
	}

	fmt.Printf("Starting server on port %s", os.Getenv("PORT"))
	log.Fatal(httpServer.ListenAndServe())
}
