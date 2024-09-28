package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"SHELLHACKS-BACKEND/auth"
	"SHELLHACKS-BACKEND/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env vars: %v", err)
	}

	// Initialize the Authenticator
	auth, err := auth.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	// Create the Gin router
	rtr := routes.New(auth)

	// Start the server
	log.Print("Server starting on http://localhost:3000/")
	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
