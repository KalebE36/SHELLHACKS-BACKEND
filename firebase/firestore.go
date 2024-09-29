package firebase

import (
	"context"
	"fmt"
	"log"
	"strings"

	firebase "firebase.google.com/go"
)

// CreateUser generates a random userID, extracts the username from the email, and adds a new user document to Firestore.
// The photoURL parameter is optional.
func CreateUser(app *firebase.App, userID string, email string, photoURL *string) error {
	// Initialize Firestore client
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
		return err
	}
	defer client.Close()

	// Extract username from email (everything before the '@')
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return fmt.Errorf("invalid email address: %s", email)
	}
	username := email[:atIndex]

	// Prepare the userData map
	userData := map[string]interface{}{
		"username": username,
		"email":    email,
	}

	// If a photoURL is provided, add it to the userData map
	if photoURL != nil && *photoURL != "" {
		userData["photo_url"] = *photoURL
	} else {
		userData["photo_url"] = "https://postimg.cc/4nY9fSWj"
	}

	// Create a new document in the "users" collection with the generated user data
	_, err = client.Collection("users").Doc(userID).Set(ctx, userData)
	if err != nil {
		log.Fatalf("Failed creating user: %v", err)
		return err
	}

	fmt.Printf("User created successfully with userID: %s and username: %s\n", userID, username)
	return nil
}

// AddSubCollection adds a document to a sub-collection within a user document.
func AddSubCollection(app *firebase.App, userID, collectionName string, data map[string]interface{}) error {
	// Initialize Firestore client.
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
		return err
	}
	defer client.Close()

	// Add a document to a sub-collection under the user document
	subCollectionRef := client.Collection("users").Doc(userID).Collection(collectionName)

	// Add the data to the sub-collection
	_, _, err = subCollectionRef.Add(ctx, data)
	if err != nil {
		log.Fatalf("Failed adding to sub-collection: %v", err)
		return err
	}

	fmt.Println("Data added to sub-collection successfully!")
	return nil
}
