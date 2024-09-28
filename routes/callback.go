package routes

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"SHELLHACKS-BACKEND/auth"
	"SHELLHACKS-BACKEND/database"
	"SHELLHACKS-BACKEND/models"

	"cloud.google.com/go/firestore"
	"golang.org/x/oauth2"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Modularized CallbackHandler
func CallbackHandler(auth *auth.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Validate state parameter
		if !validateState(ctx, session) {
			return
		}

		// Exchange authorization code for token
		token, err := exchangeCodeForToken(ctx, auth)
		if err != nil {
			return
		}

		// Verify the ID token and extract user profile
		profile, err := verifyIDTokenAndExtractProfile(ctx, auth, token)
		if err != nil {
			return
		}

		// Save profile and access token in session
		if !saveSession(ctx, session, profile, token) {
			return
		}

		// Initialize Firestore client
		fsClient, err := initializeFirestoreClient(ctx)
		if err != nil {
			return
		}
		defer fsClient.Close()

		// Handle user creation or update in Firestore
		if !handleFirestoreUser(ctx, fsClient, profile) {
			return
		}

		// Redirect to user profile page
		ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:4321/")
	}
}

// Modularized helper functions

// Validate the state parameter to prevent CSRF attacks
func validateState(ctx *gin.Context, session sessions.Session) bool {
	state := ctx.Query("state")
	if state != session.Get("state") {
		log.Println("State mismatch. Potential CSRF attack.")
		ctx.String(http.StatusBadRequest, "Invalid state parameter.")
		return false
	}
	return true
}

// Exchange authorization code for token
func exchangeCodeForToken(ctx *gin.Context, auth *auth.Authenticator) (*oauth2.Token, error) {
	code := ctx.Query("code")
	token, err := auth.Exchange(ctx.Request.Context(), code)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		ctx.String(http.StatusUnauthorized, "Failed to exchange authorization code.")
		return nil, err
	}
	return token, nil
}

// Verify the ID token and extract the user profile
func verifyIDTokenAndExtractProfile(ctx *gin.Context, auth *auth.Authenticator, token *oauth2.Token) (map[string]interface{}, error) {
	idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
	if err != nil {
		log.Printf("ID token verification failed: %v", err)
		ctx.String(http.StatusInternalServerError, "Failed to verify ID token.")
		return nil, err
	}

	// Extract user profile
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		log.Printf("Failed to parse ID token claims: %v", err)
		ctx.String(http.StatusInternalServerError, "Failed to parse claims.")
		return nil, err
	}

	return profile, nil
}

// Save profile and access token in session
func saveSession(ctx *gin.Context, session sessions.Session, profile map[string]interface{}, token *oauth2.Token) bool {
	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)
	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %v", err)
		ctx.String(http.StatusInternalServerError, "Failed to save session.")
		return false
	}
	return true
}

// Initialize Firestore client
func initializeFirestoreClient(ctx *gin.Context) (*firestore.Client, error) {
	fsClient, err := database.InitializeFirestore()
	if err != nil {
		log.Printf("Failed to initialize Firestore: %v", err)
		ctx.String(http.StatusInternalServerError, "Failed to initialize Firestore.")
		return nil, err
	}
	return fsClient, nil
}

// Handle Firestore user creation or update
func handleFirestoreUser(ctx *gin.Context, fsClient *firestore.Client, profile map[string]interface{}) bool {
	userID := profile["sub"].(string)
	userDoc := fsClient.Collection("users").Doc(userID)
	doc, err := userDoc.Get(ctx.Request.Context())
	if err != nil && !doc.Exists() {
		// User does not exist, create a new user
		if !createNewFirestoreUser(ctx.Request.Context(), userDoc, profile) {
			return false
		}
	} else {
		log.Printf("User already exists in Firestore: %v", profile["email"])
	}
	return true
}

// Create a new user in Firestore
func createNewFirestoreUser(ctx context.Context, userDoc *firestore.DocumentRef, profile map[string]interface{}) bool {
	// Safely extract profile data
	email, ok := profile["email"].(string)
	if !ok || email == "" {
		log.Println("Email not found in profile")
		return false
	}

	picture, _ := profile["picture"].(string)   // Safe type assertion, default to empty string
	nickname, _ := profile["nickname"].(string) // Same here, using safe type assertion

	// Create the user model
	user := models.User{
		ID:       profile["sub"].(string), // Assuming 'sub' is always present and valid
		Email:    email,
		Picture:  picture,
		Username: nickname, // Using nickname if available
	}

	// Add the user to Firestore
	_, err := userDoc.Set(ctx, user)
	if err != nil {
		log.Printf("Failed to add user to Firestore: %v", err)
		return false
	}

	log.Printf("User added to Firestore: %v", user.Email)
	return true
}

// HandleCardCollections handles adding a card to the appropriate collection
func HandleCardCollections(ctx context.Context, userDoc *firestore.DocumentRef, cardType string, word string, number int) bool {
	log.Printf("Handling card collections for %s...\n", cardType)

	if !addCardToCollection(ctx, userDoc, cardType, word, number) {
		log.Printf("Failed to add card to collection: %s\n", cardType)
		return false
	}

	log.Printf("Successfully handled card collections for %s\n", cardType)
	return true
}

// TestHandleCardCollections handles adding a card to the appropriate collection for testing
func TestHandleCardCollections(ctx context.Context, userDoc *firestore.DocumentRef, cardType string, word string, number int) bool {
	return HandleCardCollections(nil, userDoc, cardType, word, number) // Use nil or pass context for testing
}

// Add a card to the appropriate collection when the user learns or encounters a card for the first time
func addCardToCollection(ctx context.Context, userDoc *firestore.DocumentRef, collectionType string, word string, number int) bool {
	// Get the correct collection reference
	collection := userDoc.Collection("flashcards").Doc(collectionType).Collection(collectionType)

	// If we're adding to unknowncards, check the size of the collection first
	if collectionType == "unknowncards" {
		// Get the documents in the unknowncards collection
		docs, err := collection.Documents(ctx).GetAll()
		if err != nil {
			log.Printf("Failed to retrieve documents in %v collection: %v", collection.Path, err)
			return false
		}

		// If the collection has more than 30 documents, move the new card to overflowcards
		if len(docs) >= 30 {
			log.Printf("unknowncards collection has reached the limit. Adding card to overflowcards.")
			collectionType = "overflowcards" // Change the collection to overflowcards
			collection = userDoc.Collection("flashcards").Doc(collectionType).Collection(collectionType)
		}
	}

	// Generate a random document ID for the new card
	newDocID, err := generateRandomID()
	if err != nil {
		log.Printf("Failed to generate random ID for document in %v collection: %v", collection.Path, err)
		return false
	}

	log.Printf("Attempting to create document in %v collection...\n", collection.Path)

	// Add the card with the generated ID
	_, err = collection.Doc(newDocID).Set(ctx, map[string]interface{}{
		"word":      word,
		"number":    number,
		"createdAt": firestore.ServerTimestamp,
	})

	if err != nil {
		log.Printf("Failed to create document in %v collection: %v", collection.Path, err)
		return false
	}

	log.Printf("Document %s created in %v collection", newDocID, collection.Path)
	return true
}

// Generate a random document ID
func generateRandomID() (string, error) {
	// Create a random 16-byte ID and encode it as a hex string
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
