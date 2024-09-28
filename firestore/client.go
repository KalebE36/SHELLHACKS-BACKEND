package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"        // Firestore client package
	"firebase.google.com/go"               // Firebase Admin SDK package
	"google.golang.org/api/option"         // For providing service account key
)

// InitializeFirestore initializes the Firestore client
func InitializeFirestore() (*firestore.Client, error) {
	// Set the path to your Firebase service account key
	opt := option.WithCredentialsFile("shellhacks-f1d7c-firebase-adminsdk-f9j4e-4d541506af.json")

	// Initialize the Firebase App
	app, err := firebase.NewApp(context.Background(), nil, opt)
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
