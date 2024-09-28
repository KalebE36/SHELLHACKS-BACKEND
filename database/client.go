package database

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

// Global variable to hold the Firebase app instance
var firebaseApp *firebase.App
var errFirebaseNotInitialized = errors.New("firebase not initialized")

// InitializeFirebase initializes the Firebase app with service account credentials.
func InitializeFirebase() error {
	ctx := context.Background()

	// Path to your service account key JSON file
	sa := option.WithCredentialsFile("database/shellhacks-go-firebase-adminsdk-kud5n-4be4c82471.json")

	// Initialize the Firebase app
	var err error
	firebaseApp, err = firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Printf("Failed to initialize Firebase app: %v", err)
		return err
	}

	log.Println("Firebase app initialized successfully")
	return nil
}

// InitializeFirestoreClient initializes and returns the Firestore client.
func InitializeFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	// Ensure Firebase is initialized
	if firebaseApp == nil {
		return nil, errFirebaseNotInitialized
	}

	// Initialize Firestore client
	projectID := "your-project-id" // Replace with your actual Firebase project ID
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("Failed to initialize Firestore client: %v", err)
		return nil, err
	}

	return client, nil
}

// GetAuthClient initializes and returns the Firebase Auth client.
func GetAuthClient(ctx context.Context) (*auth.Client, error) {
	// Ensure Firebase is initialized
	if firebaseApp == nil {
		return nil, errFirebaseNotInitialized
	}

	// Initialize Firebase Auth client
	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Printf("Failed to initialize Firebase Auth client: %v", err)
		return nil, err
	}

	return authClient, nil
}
