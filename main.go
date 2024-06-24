package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/deadpyxel/curator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import postgres drive and use side effects
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		log.Fatal("PORT is not defined")
	}

	connString := os.Getenv("CONN_STRING")
	if connString == "" {
		log.Fatal("CONN_STRING is not defined")
	}

	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	// Heathcheck
	mux.HandleFunc("GET /v1/healthz", handlerLiveness)
	mux.HandleFunc("GET /v1/err", handlerErrorTest)
	// Users
	mux.HandleFunc("POST /v1/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", apiCfg.authMiddleware(apiCfg.handlerGetUser))
	// Feeds
	mux.HandleFunc("POST /v1/feeds", apiCfg.authMiddleware(apiCfg.handlerCreateFeed))
	mux.HandleFunc("GET /v1/feeds", apiCfg.handlerGetFeeds)
	mux.HandleFunc("POST /v1/feed_follows", apiCfg.authMiddleware(apiCfg.handlerCreateFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.authMiddleware(apiCfg.handlerGetFeedFollows))
	logMux := logMiddleware(mux)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: logMux,
	}

	fmt.Printf("Starting server on port %s\n\n", serverPort)
	log.Fatal(httpServer.ListenAndServe())
}
