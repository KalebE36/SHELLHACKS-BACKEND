package routes

import (
	"log"
	"net/http"

	"SHELLHACKS-BACKEND/database"
	"SHELLHACKS-BACKEND/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CallbackHandler handles the callback after Firebase authentication
func CallbackHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Validate the state parameter to prevent CSRF attacks
		state := ctx.Query("state")
		if state != session.Get("state") {
			log.Println("State mismatch. Potential CSRF attack.")
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Get the ID token from the request
		idToken := ctx.Query("id_token")
		if idToken == "" {
			log.Println("ID token is missing")
			ctx.String(http.StatusUnauthorized, "Missing ID token")
			return
		}

		// Initialize Firebase Auth client
		authClient, err := database.GetAuthClient(ctx.Request.Context())
		if err != nil {
			log.Printf("Failed to initialize Firebase Auth client: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to initialize Firebase Auth")
			return
		}

		// Verify the ID token
		token, err := authClient.VerifyIDToken(ctx.Request.Context(), idToken)
		if err != nil {
			log.Printf("Failed to verify ID token: %v", err)
			ctx.String(http.StatusUnauthorized, "Failed to verify ID token")
			return
		}

		// Extract user information from the token
		uid := token.UID

		// Save user info in the session
		session.Set("user_id", uid)
		if err := session.Save(); err != nil {
			log.Printf("Failed to save session: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to save session")
			return
		}

		// Initialize Firestore client
		fsClient, err := database.InitializeFirestoreClient(ctx.Request.Context())
		if err != nil {
			log.Printf("Failed to initialize Firestore client: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to initialize Firestore")
			return
		}
		defer fsClient.Close()

		// Check if the user already exists in Firestore
		userDoc := fsClient.Collection("users").Doc(uid)
		doc, err := userDoc.Get(ctx.Request.Context())
		if err != nil && !doc.Exists() {
			// User does not exist, create a new user
			user := models.User{
				ID: uid,
			}

			// Add the user to Firestore
			_, err = userDoc.Set(ctx.Request.Context(), user)
			if err != nil {
				log.Printf("Failed to add user to Firestore: %v", err)
				ctx.String(http.StatusInternalServerError, "Failed to add user to Firestore")
				return
			}
			log.Printf("User added to Firestore: %v", uid)
		} else {
			log.Printf("User already exists in Firestore: %v", uid)
		}

		// Redirect to the frontend page (replace with your frontend URL)
		ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:4321/")
	}
}
