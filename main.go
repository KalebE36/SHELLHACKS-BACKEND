package main

import (
	"SHELLHACKS-BACKEND/firebase" // Import your custom firebase package
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize Firebase app
	firebaseApp, err := firebase.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	photo_url := "https://lh3.googleusercontent.com/a/ACg8ocLJqvbrBQPZGofc6vh9Up4hZ1jtl01eS-JFM_iSydyHna5jB_bhVQ=s96-c"
	// Create a new user in Firestore
	err = firebase.CreateUser(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e70", "example@example.com", &photo_url)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}

	err = firebase.CreateUser(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e71", "bob@example.com", nil)
	// Add a sub-collection to the user document
	data := map[string]interface{}{
		"word": "Laptop",
		"int":  1200,
	}
	err = firebase.AddSubCollection(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e70", "known_card", data)
	if err != nil {
		log.Fatalf("Error adding to sub-collection: %v", err)
	}

	data = map[string]interface{}{
		"word": "quack",
		"int":  120,
	}

	err = firebase.AddSubCollection(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e70", "known_card", data)

	// Initialize the router from routes.go (as before)
	router := http.NewServeMux()
	log.Fatal(http.ListenAndServe(":3000", router))
}
