package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"SHELLHACKS-BACKEND/auth"
	"SHELLHACKS-BACKEND/routes"

	"github.com/joho/godotenv"
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

	log.Println("GEMINI_API_KEY:", os.Getenv("GEMINI_API_KEY"))

	// Run the test before starting the server
	// testAddCardsToFirestore()

	// Initialize the Authenticator
	authenticator, err := auth.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	gob.Register(map[string]interface{}{})

	// Create the Gin router
	router := routes.New(authenticator)

	// Start the server
	if err := http.ListenAndServe("0.0.0.0:3000", router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// // testAddCardsToFirestore adds cards directly to Firestore for testing purposes
// func testAddCardsToFirestore() {
// 	log.Println("Running testAddCardsToFirestore...")

// 	// Initialize Firestore client
// 	fsClient, err := database.InitializeFirestore()
// 	if err != nil {
// 		log.Fatalf("Failed to initialize Firestore: %v", err)
// 	}
// 	defer fsClient.Close()

// 	// Simulate a test user ID
// 	userID := "google-oauth2|116640021460997271224"
// 	userDoc := fsClient.Collection("users").Doc(userID)

// 	log.Println("Testing adding cards...")

// 	// Test adding cards to different collections
// 	testAddCard(fsClient, userDoc, "knowncards", "apple", 1)
// 	testAddCard(fsClient, userDoc, "unknowncards", "banana1", 2)
// 	testAddCard(fsClient, userDoc, "overflowcards", "cherry", 3)

// 	log.Println("Finished adding cards.")
// }

// // testAddCard simulates adding a card to a specific collection
// func testAddCard(fsClient *firestore.Client, userDoc *firestore.DocumentRef, cardType string, word string, number int) {
// 	log.Printf("Attempting to add card to %s: %s (%d)\n", cardType, word, number)

// 	success := routes.HandleCardCollections(context.Background(), userDoc, cardType, word, number)
// 	if success {
// 		log.Printf("Successfully added card to %s: %s (%d)\n", cardType, word, number)
// 	} else {
// 		log.Printf("Failed to add card to %s\n", cardType)
// 	}
// }
