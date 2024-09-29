package main

import (
	"SHELLHACKS-BACKEND/routes"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize the router from router.go
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router := routes.InitializeRouter()
	

	// Start the server on port 3000
	if err := http.ListenAndServe("0.0.0.0:3000", router); err != nil {
		log.Fatal(err)
	}
}
