package main

import (
	"encoding/gob"
	"log"
	"net/http"

	"SHELLHACKS-BACKEND/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Register map[string]interface{} type with gob
	gob.Register(map[string]interface{}{})
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env vars: %v", err)
	}

	// Log the loaded callback URL for debugging purposes

	// Initialize the Authenticator

	// Create a Gin router
	router := gin.Default()

	// Define routes using your existing routes package
	router.GET("/callback", routes.CallbackHandler())

	// Start the server on port 3000
	if err := http.ListenAndServe("0.0.0.0:3000", router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
