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

// AddSubCollectionToLanguage checks if the "unknown_cards" sub-collection exists and has over 30 items,
// then switches to "overflow_cards" if needed.
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

	// Use a query to count the number of documents in the "unknown_cards" sub-collection
	unknownCardsCollectionRef := languageDocRef.Collection("unknown_cards")
	docsCount, err := unknownCardsCollectionRef.Limit(31).Documents(ctx).GetAll() // Fetch a maximum of 31 documents to check if the limit is exceeded
	if err != nil {
		log.Fatalf("Failed to retrieve unknown_cards sub-collection: %v", err)
		return err
	}

	fmt.Printf("Found %d documents in 'unknown_cards'\n", len(docsCount)) // Debug print to show the document count

	// If there are more than or equal to 30 documents in "unknown_cards", switch to "overflow_cards"
	if len(docsCount) >= 30 {
		subCollectionName = "overflow_cards"
		fmt.Println("Switching to overflow_cards as unknown_cards has reached the limit of 30 items")
	}

	// Reference to the sub-collection within the language document
	subCollectionRef := languageDocRef.Collection(subCollectionName)

	// Add a document to the sub-collection with the provided data
	_, _, err = subCollectionRef.Add(ctx, data)
	if err != nil {
		log.Fatalf("Failed to add data to sub-collection: %v", err)
		return err
	}

	fmt.Printf("Data added to sub-collection %s in language %s with provided data!\n", subCollectionName, language)
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
		"expertise":       0,
	}, firestore.MergeAll) // MergeAll ensures that only the "integer" field is added or updated
	if err != nil {
		log.Fatalf("Failed to add integer field to spanish document: %v", err)
		return err
	}

	fmt.Println("Language and spanish document created successfully with integer field set to 0!")
	return nil
}

// UpdateLanguageField adds the passed value to the current field value in the language document.
func UpdateLanguageField(app *firebase.App, userID string, language string, field string, incrementValue int) error {
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

	// Retrieve the current value of the specified field
	docSnapshot, err := languageDocRef.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to retrieve language document: %v", err)
		return err
	}

	// Get the current value of the field, assuming the field is of type int
	currentValue, ok := docSnapshot.Data()[field].(int64) // Firestore stores integers as int64
	if !ok {
		return fmt.Errorf("field %s is not of type int", field)
	}

	// Add the passed value to the current value
	newValue := currentValue + int64(incrementValue)

	// Update the field with the new value
	_, err = languageDocRef.Update(ctx, []firestore.Update{
		{Path: field, Value: newValue},
	})
	if err != nil {
		log.Fatalf("Failed to update field in language document: %v", err)
		return err
	}

	fmt.Printf("Field %s in language %s updated successfully! New value: %d\n", field, language, newValue)
	return nil
}

func createFirestore() *firestore.Client {
	projectID := "shellhacks-f1d7c"
	client, err := firestore.NewClient(context.Background(), projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}
