package database

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"   // Firestore client package
	firebase "firebase.google.com/go" // Firebase Admin SDK package
	"google.golang.org/api/option"    // For providing service account key
)

// InitializeFirestore initializes the Firestore client
func InitializeFirestore() (*firestore.Client, error) {
	// Set the path to your Firebase service account key
	opt := option.WithCredentialsFile("database/shellhacks-go-firebase-adminsdk-kud5n-4be4c82471.json")

	// Initialize the Firebase App
	config := &firebase.Config{
		ProjectID: "shellhacks-go", // Replace with your Firebase project ID
	}

	// Initialize Firebase App with explicit project ID
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	// Initialize Firestore client
	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing Firestore client: %v", err)
	}

	return client, nil
}
