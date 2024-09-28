package routes

import (
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"SHELLHACKS-BACKEND/auth"
)

func New(auth *auth.Authenticator) *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies, we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	// Initialize cookie store and session management
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// Serve static files and HTML templates

	// Define routes

	router.GET("/login", LoginHandler(auth))
	router.GET("/callback", CallbackHandler(auth))
	router.GET("/logout", LogoutHandler)

	return router
}
