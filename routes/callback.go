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
		if ctx.Query("state") != session.Get("state") {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange the authorization code for a token
		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Failed to exchange the authorization code.")
			return
		}

		// Verify the ID token
		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to verify ID token.")
			return
		}

		// Extract and save user profile in session
		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		fsClient, err := firestore.InitializeFirestore()
		if err != nil {
			log.Printf("Failed to initialize Firestore: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to initialize Firestore.")
			return
		}
		defer fsClient.Close()

		// Check if user already exists in Firestore
		userID := profile["sub"].(string) // 'sub' is the unique identifier for the user
		userDoc := fsClient.Collection("users").Doc(userID)
		doc, err := userDoc.Get(ctx.Request.Context())
		if err != nil && !doc.Exists() {
			// User does not exist, so create a new User instance
			user := models.User{
				ID:        userID,
				FirstName: profile["given_name"].(string),
				LastName:  profile["family_name"].(string),
				Email:     profile["email"].(string),
				Picture:   profile["picture"].(string),
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
			log.Printf("User already exists in Firestore: %v", profile["email"])
		}

		// Redirect to the user page
		ctx.Redirect(http.StatusTemporaryRedirect, "/user")
	}
}
