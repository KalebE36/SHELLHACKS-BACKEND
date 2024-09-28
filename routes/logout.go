package routes

import (
	"net/http"
	"net/url"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler for logging out the user.
func LogoutHandler(ctx *gin.Context) {
	// Clear the session
	session := sessions.Default(ctx)
	session.Clear()       // Remove all session data
	err := session.Save() // Save the empty session to remove it
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to clear session")
		return
	}

	// Construct the Auth0 logout URL
	logoutURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Determine the return URL (after logout)
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Add the necessary query parameters
	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))

	logoutURL.RawQuery = parameters.Encode()

	// Redirect to the Auth0 logout URL
	ctx.Redirect(http.StatusTemporaryRedirect, logoutURL.String())
}
