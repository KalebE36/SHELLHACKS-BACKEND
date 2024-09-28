package cards

import (
	"SHELLHACKS-BACKEND/auth"
	"SHELLHACKS-BACKEND/database"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ReturnCardsHandler(auth *auth.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Get the user's profile from the session
		profile := session.Get("profile")
		if profile == nil {
			log.Println("No profile found in session")
			ctx.String(http.StatusUnauthorized, "User not authenticated")
			return
		}

		// Type assertion to convert profile back to a map
		profileMap, ok := profile.(map[string]interface{})
		if !ok {
			log.Println("Failed to convert profile to map")
			ctx.String(http.StatusInternalServerError, "Failed to retrieve profile")
			return
		}

		// Extract the user's unique ID (e.g., 'sub')
		userID, ok := profileMap["sub"].(string)
		if !ok {
			log.Println("Failed to extract user ID from profile")
			ctx.String(http.StatusInternalServerError, "Failed to retrieve user ID")
			return
		}

		// Initialize Firestore client
		fsClient, err := database.InitializeFirestore()
		if err != nil {
			log.Printf("Failed to initialize Firestore: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to initialize Firestore")
			return
		}
		defer fsClient.Close()

		// Define paths to knowncards and unknowncards collections
		knownCardsCollection := fsClient.Collection("users").Doc(userID).Collection("knowncards")
		unknownCardsCollection := fsClient.Collection("users").Doc(userID).Collection("unknowncards")

		// Fetch documents from knowncards
		knownDocs, err := knownCardsCollection.Documents(ctx.Request.Context()).GetAll()
		if err != nil {
			log.Printf("Failed to fetch known cards: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to fetch known cards")
			return
		}

		// Fetch documents from unknowncards
		unknownDocs, err := unknownCardsCollection.Documents(ctx.Request.Context()).GetAll()
		if err != nil {
			log.Printf("Failed to fetch unknown cards: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to fetch unknown cards")
			return
		}

		// Prepare slices to store the cards
		var knownCards []map[string]interface{}
		var unknownCards []map[string]interface{}

		// Loop through knowncards documents and add them to the slice
		for _, doc := range knownDocs {
			knownCards = append(knownCards, doc.Data())
		}

		// Loop through unknowncards documents and add them to the slice
		for _, doc := range unknownDocs {
			unknownCards = append(unknownCards, doc.Data())
		}

		// Return both known and unknown cards as a JSON response
		ctx.JSON(http.StatusOK, gin.H{
			"knownCards":   knownCards,
			"unknownCards": unknownCards,
		})
	}
}
