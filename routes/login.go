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
		state, err := helpers.GenerateState()

		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		// Store the state in the session
		session := sessions.Default(ctx)
		session.Set("state", state)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to the Auth0 login page
		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
	}
}
