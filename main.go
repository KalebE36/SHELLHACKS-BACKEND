package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"SHELLHACKS-BACKEND/auth"
	"SHELLHACKS-BACKEND/routes"
)

func init() {
	// Register map[string]interface{} type with gob
	gob.Register(map[string]interface{}{})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env vars: %v", err)
	}

	log.Println("Loaded callback URL:", os.Getenv("AUTH0_CALLBACK_URL"))

	// Initialize the Authenticator
	auth, err := auth.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	gob.Register(map[string]interface{}{})
	// Create the Gin router
	router := routes.New(auth)

	// Start the server
	if err := http.ListenAndServe("0.0.0.0:3000", router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
