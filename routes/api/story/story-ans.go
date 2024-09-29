package story

import (
	"net/http"

	"SHELLHACKS-BACKEND/helpers"

	"github.com/gin-gonic/gin"
)

func HandleStoryResponse(ctx *gin.Context) {
	var requestBody struct {
		UserResponse string `json:"user_response"`
		CurrentStory string `json:"current_story"`
	}

	// Bind the JSON from the request body
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input."})
		return
	}

	// Generate the next story segment using the user's response
	newSegment, err := helpers.GenerateStoryHandler(requestBody.CurrentStory, requestBody.UserResponse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate story segment."})
		return
	}

	// Return the new segment of the story
	ctx.JSON(http.StatusOK, gin.H{"new_segment": newSegment})
}
