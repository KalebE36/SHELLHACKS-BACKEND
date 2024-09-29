package user

import (
	"log"
	"net/http"

	"fmt"

	"SHELLHACKS-BACKEND/firebase"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(ctx *gin.Context) {
	var incomingData struct {
		UID      string `json:"uid"`
		Email    string `json:"email"`
		PhotoURL string `json:"photoURL"`
	}

	// Unmarshal the JSON into the struct
	// Bind the JSON from the request body
	if err := ctx.ShouldBindJSON(&incomingData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input."})
		return
	}

	// Initialize Firebase app
	firebaseApp, err := firebase.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Print the extracted fields
	fmt.Println("UID:", incomingData.UID)
	fmt.Println("Email:", incomingData.Email)
	fmt.Println("PhotoURL:", incomingData.PhotoURL)

	photo_url := &incomingData.PhotoURL
	firebase.CreateUser(firebaseApp, incomingData.UID, incomingData.Email, photo_url)
}
