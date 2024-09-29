package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// InitializeApp initializes the Firebase app with the service account credentials.
func InitializeApp() (*firebase.App, error) {
	// Use the service account key JSON file for authentication.
	sa := option.WithCredentialsFile("firebase/shellhacks-go-firebase-adminsdk-kud5n-4be4c82471.json")

	// Set the project ID explicitly
	conf := &firebase.Config{
		ProjectID: "shellhacks-f1d7c", // Replace with your Firebase project ID
	}

	// Initialize the Firebase app with the configuration
	app, err := firebase.NewApp(context.Background(), conf, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return nil, err
	}

	return app, nil
}
