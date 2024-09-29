package firebase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/firestore"
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

// AddSubCollectionToLanguage adds a sub-collection to a specific language document (e.g., "spanish") inside the "language" collection.
func AddSubCollectionToLanguage(app *firebase.App, userID string, language string, subCollectionName string, data map[string]interface{}) error {
	// Initialize Firestore client
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
		return err
	}
	defer client.Close()

	// Reference to the specific language document under the "language" collection
	languageDocRef := client.Collection("users").Doc(userID).Collection("language").Doc(language)

	// Reference to the sub-collection within the language document
	subCollectionRef := languageDocRef.Collection(subCollectionName)

	// Add a document to the sub-collection with the provided data
	_, _, err = subCollectionRef.Add(ctx, data)
	if err != nil {
		log.Fatalf("Failed to add data to sub-collection: %v", err)
		return err
	}

	fmt.Printf("Sub-collection %s added to language %s with provided data!\n", subCollectionName, language)
	return nil
}

// AddLanguage creates a "language" collection under a user and a "spanish" document with a direct "integer" field.
func AddLanguage(app *firebase.App, userID string, language string) error {
	// Initialize Firestore client
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
		return err
	}
	defer client.Close()

	// Reference to the "language" collection under the user document
	languageCollectionRef := client.Collection("users").Doc(userID).Collection("language")

	// Reference to the "spanish" document under the "language" collection
	spanishDocRef := languageCollectionRef.Doc(language)

	// Add the "integer" field directly to the "spanish" document
	_, err = spanishDocRef.Set(ctx, map[string]interface{}{
		"total_cards":     0,
		"learned_cards":   0,
		"paragraphs_read": 0,
		"streak":          0,
	}, firestore.MergeAll) // MergeAll ensures that only the "integer" field is added or updated
	if err != nil {
		log.Fatalf("Failed to add integer field to spanish document: %v", err)
		return err
	}

	fmt.Println("Language and spanish document created successfully with integer field set to 0!")
	return nil
}
