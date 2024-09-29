package user

import (
	"SHELLHACKS-BACKEND/firebase"
	"SHELLHACKS-BACKEND/models"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RetCardHandler(ctx *gin.Context) {
	var requestBody struct {
		UserId string `json:"user_id"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input."})
		return
	}

	fbClient, err := firebase.InitializeApp()
	if err != nil {
		log.Printf("Failed to initialize Firestore client: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Firestore"})
		return
	}

	fsClient, err := fbClient.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer fsClient.Close()

	cardsRef := fsClient.Collection("users").Doc(requestBody.UserId).Collection("Spanish")
	docs, err := cardsRef.Documents(context.Background()).GetAll()
	if err != nil {
		log.Printf("Failed to retrieve cards for user %s: %v", requestBody.UserId, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cards from Firestore"})
		return
	}

	// Create an array to hold the retrieved cards
	var cards []models.Flashcard

	// Loop through the documents and decode them into the Flashcard struct
	for _, doc := range docs {
		var card models.Flashcard
		if err := doc.DataTo(&card); err != nil {
			log.Printf("Failed to decode card data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode card data"})
			return
		}
		cards = append(cards, card)
	}

	// Return the array of cards as a JSON response
	ctx.JSON(http.StatusOK, gin.H{"cards": cards})
}
