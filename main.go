package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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
