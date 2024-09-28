package routes

import (
	"log"
	"net/http"

	"SHELLHACKS-BACKEND/auth"
	"SHELLHACKS-BACKEND/firestore"
	"SHELLHACKS-BACKEND/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler for the callback after login.
func CallbackHandler(auth *auth.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		state := ctx.Query("state")
		if state != session.Get("state") {
			log.Println("State mismatch. Potential CSRF attack.")
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange the authorization code for a token
		code := ctx.Query("code")
		token, err := auth.Exchange(ctx.Request.Context(), code)
		if err != nil {
			log.Printf("Token exchange failed: %v", err)
			ctx.String(http.StatusUnauthorized, "Failed to exchange the authorization code.")
			return
		}

		// Verify the ID token
		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			log.Printf("ID token verification failed: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to verify ID token.")
			return
		}

		// Extract user profile and store it in session
		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			log.Printf("Failed to parse ID token claims: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to parse claims.")
			return
		}

		// Safely extract profile data with type assertions
		firstName, _ := profile["given_name"].(string)
		lastName, _ := profile["family_name"].(string)
		email, _ := profile["email"].(string)
		picture, _ := profile["picture"].(string)

		// Save access token and profile to the session
		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			log.Printf("Failed to save session: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to save session.")
			return
		}

		// Initialize Firestore client
		fsClient, err := firestore.InitializeFirestore()
		if err != nil {
			log.Printf("Failed to initialize Firestore: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to initialize Firestore.")
			return
		}
		defer fsClient.Close()

		// Check if user already exists in Firestore
		userID := profile["sub"].(string) // 'sub' is a unique identifier for the user
		userDoc := fsClient.Collection("users").Doc(userID)
		doc, err := userDoc.Get(ctx.Request.Context())
		if err != nil && !doc.Exists() {
			// User does not exist, create a new User instance
			user := models.User{
				ID:        userID,
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Picture:   picture,
			}

			// Add the user to Firestore
			_, err := userDoc.Set(ctx.Request.Context(), user)
			if err != nil {
				log.Printf("Failed to add user to Firestore: %v", err)
				ctx.String(http.StatusInternalServerError, "Failed to add user to Firestore.")
				return
			}
			log.Printf("User added to Firestore: %v", user.Email)
		} else {
			log.Printf("User already exists in Firestore: %v", email)
		}

		// Redirect to user profile page
		ctx.Redirect(http.StatusTemporaryRedirect, "/user")
	}
}
