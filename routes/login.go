package routes

import (
	"SHELLHACKS-BACKEND/auth"
	"SHELLHACKS-BACKEND/helpers"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"net/http"
)

// Handler for the login process.
func LoginHandler(auth *auth.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Generate a state value to prevent CSRF attacks
		state, err := helpers.GenerateState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate state")
			return
		}

		// Save state to session
		session.Set("state", state)
		session.Save()

		// Redirect to Auth0 for login
		authURL := auth.AuthCodeURL(state, oauth2.AccessTypeOffline)
		log.Printf("Redirecting to Auth0 with URL: %s", authURL)
		ctx.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}
