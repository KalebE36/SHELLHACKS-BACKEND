package main

import (
	"SHELLHACKS-BACKEND/routes"
	"log"
	"net/http"
)

func main() {
	// Initialize the router from router.go
	router := routes.InitializeRouter()

	// Start the server on port 3000
	if err := http.ListenAndServe("0.0.0.0:3000", router); err != nil {
		log.Fatal(err)
	}
}
