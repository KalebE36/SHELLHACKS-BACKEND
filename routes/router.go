package routes

import (
	"SHELLHACKS-BACKEND/helpers"
	"SHELLHACKS-BACKEND/routes/api"

	"github.com/gorilla/mux"
)

// InitializeRouter sets up the routes for the application
func InitializeRouter() *mux.Router {
	// Create a new router
	router := mux.NewRouter()

	// Define the routes and map them to handler functions
	router.HandleFunc("/api/generate-content", helpers.ConvertGinToMux(api.GenerateParagraphsHandler)).Methods("POST")

	// Add more routes as needed
	// router.HandleFunc("/other", OtherHandler).Methods("GET")

	return router
}
