package user

import (
	"context"
	"log"
	"net/http"

	"SHELLHACKS-BACKEND/firebase"
	"SHELLHACKS-BACKEND/models"

	"cloud.google.com/go/firestore"

	"github.com/gin-gonic/gin"
)

func SpacedRepetition(card *models.Flashcard) {
	*card.Proficiency = *card.Proficiency * 2 // Double the interval if the answer was correct
}

func checkIfCardExists(client *firestore.Client, userID string, cardWord string) (bool, error) {
	// Get reference to the user's flashcards collection
	cardsRef := client.Collection("users").Doc(userID).Collection("Spanish")

	// Reference the document by cardWord (document ID)
	cardDoc := cardsRef.Doc(cardWord)

	// Attempt to get the document snapshot
	docSnap, err := cardDoc.Get(context.Background())
	if err != nil {
		// Return the error if something goes wrong other than document not found
		return false, err
	}

	// Check if the document exists
	return docSnap.Exists(), nil

}

func MakeCardHandler(ctx *gin.Context) {

	var requestBody struct {
		UserId    string             `json:"user_id"`
		CardArray []models.Flashcard `json:"card_array"`
		Pass      []bool             `json:"boolean_array"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input."})
		return
	}

	// Initialize Firestore client
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

	userFlashcards := fsClient.Collection("users").Doc(requestBody.UserId).Collection("Spanish")

	for _, card := range requestBody.CardArray {
		if card.Word == "" {
			log.Printf("Card word is empty for user %s", requestBody.UserId)
			continue
		}

		_, err := userFlashcards.Doc(card.Word).Set(ctx.Request.Context(), card)
		if err != nil {
			log.Printf("Failed to add card %s to Firestore for user %s: %v", card.Word, requestBody.UserId, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add card to Firestore"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cards successfully added"})
}
