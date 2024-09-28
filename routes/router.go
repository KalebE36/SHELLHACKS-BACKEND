package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"SHELLHACKS-BACKEND/auth"
	"SHELLHACKS-BACKEND/routes/api/cards"
)

func New(auth *auth.Authenticator) *gin.Engine {
	router := gin.Default()

	// CORS configuration with credentials allowed
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:4321", "http://3.147.36.237"}, // Allow frontend origin
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // This is important for cookies/sessions
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:4321"
		},
	}

	router.Use(cors.New(corsConfig))

	// Session middleware setup
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	router.GET("/login", LoginHandler(auth))
	router.GET("/callback", CallbackHandler(auth))
	router.GET("/logout", LogoutHandler)
	router.GET("/api/cards/ret-cards", cards.ReturnCardsHandler(auth))

	return router
}
