package routes

import (
	"SHELLHACKS-BACKEND/helpers"
	"SHELLHACKS-BACKEND/routes/api"
	"SHELLHACKS-BACKEND/routes/api/story"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// InitializeRouter sets up the routes for the application
func InitializeRouter() *mux.Router {
	// Create a new router
	router := mux.NewRouter()

	// Define the routes and map them to handler functions
	router.HandleFunc("/api/onboard-gen", helpers.ConvertGinToMux(api.GenerateParagraphsHandler)).Methods("POST")
	router.HandleFunc("/api/story/story-start", helpers.ConvertGinToMux(story.StartStoryHandler)).Methods("GET")
	router.HandleFunc("/api/story/story-answer", helpers.ConvertGinToMux(story.HandleStoryResponse)).Methods("POST")

	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:4321"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	// Add more routes as needed
	// router.HandleFunc("/other", OtherHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)))
	return router
}
