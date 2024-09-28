package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
	})

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
		if !createNewFirestoreUser(ctx, userDoc, profile) {
			return false
		}

		// Handle creation of knowncards and unknowncards collections
		if !handleCardCollections(ctx, userDoc) {
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
		// Handle case where email is missing or not a string
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

// Handle knowncards and unknowncards collections
func handleCardCollections(ctx *gin.Context, userDoc *firestore.DocumentRef) bool {
	// Known Cards
	if !createCardDocument(ctx, userDoc.Collection("knowncards"), "exampleName", 42) {
		return false
	}

	// Unknown Cards
	if !createCardDocument(ctx, userDoc.Collection("unknowncards"), "someName", 24) {
		return false
	}

	return true
}

// Create a card document in the specified collection
func createCardDocument(ctx *gin.Context, collection *firestore.CollectionRef, name string, number int) bool {
	// Count the number of documents and create a new document with the next ID
	docs, err := collection.Documents(ctx.Request.Context()).GetAll()
	if err != nil {
		log.Printf("Failed to retrieve documents in %v collection: %v", collection.Path, err)
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve documents in %v collection.", collection.Path))
		return false
	}

	newDocID := fmt.Sprintf("card%d", len(docs)+1)
	_, err = collection.Doc(newDocID).Set(ctx.Request.Context(), map[string]interface{}{
		"name":   name,
		"number": number,
	})
	if err != nil {
		log.Printf("Failed to create document in %v collection: %v", collection.Path, err)
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Failed to create document in %v collection.", collection.Path))
		return false
	}

	log.Printf("Document %s created in %v collection", newDocID, collection.Path)
	return true
}
