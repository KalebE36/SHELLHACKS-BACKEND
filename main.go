package main

import (
	"SHELLHACKS-BACKEND/firebase" // Import your custom firebase package
	"log"
	"net/http"
	"strconv"

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
	err = firebase.AddSubCollectionToLanguage(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e70", "spanish", "known_card", data)
	if err != nil {
		log.Fatalf("Error adding to sub-collection: %v", err)
	}

	data = map[string]interface{}{
		"word": "quack",
		"int":  120,
	}

	for i := 1; i <= 30; i++ {
		data = map[string]interface{}{
			"word": "quack" + strconv.Itoa(i-1),
			"int":  120,
		}
		err = firebase.AddSubCollectionToLanguage(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e70", "spanish", "unknown_card", data)
	}

	data = map[string]interface{}{
		"word": "moo",
		"int":  120,
	}

	err = firebase.AddSubCollectionToLanguage(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e70", "spanish", "unknown_card", data)

	data = map[string]interface{}{
		"word": "moo2",
		"int":  120,
	}

	err = firebase.AddSubCollectionToLanguage(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e70", "spanish", "unknown_card", data)

	err = firebase.AddSubCollectionToLanguage(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e70", "spanish", "known_card", data)

	err = firebase.AddLanguage(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e70", "spanish")
	err = firebase.UpdateLanguageField(firebaseApp, "e797c4c4-2976-4248-8938-a14a656c6e70", "spanish", "learned_cards", 1)

	// Initialize the router from routes.go (as before)
	router := http.NewServeMux()
	log.Fatal(http.ListenAndServe(":3000", router))
}
