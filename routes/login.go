package routes

import (
	"SHELLHACKS-BACKEND/auth"
	"SHELLHACKS-BACKEND/helpers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

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
		authURL := auth.AuthCodeURL(state)
		ctx.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}
