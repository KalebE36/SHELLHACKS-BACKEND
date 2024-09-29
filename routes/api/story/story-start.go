package story

import (
	"net/http"

	"SHELLHACKS-BACKEND/helpers"

	"github.com/gin-gonic/gin"
)

// StartStoryHandler handles the request to start the story
func StartStoryHandler(ctx *gin.Context) {
	initialPrompt := "Come up with a random story that is both exciting and random, also dont talk in the first person. This story should be different everytime and could be a story about literally anything. Imagine youre directing a decision based episode or movie and the user has options which have consequences or something similar. Just give like a couple of sentences and give 3 options for how they should respond. You dont need to add anything after that." // Initial story prompt

	// Call Gemini API to generate the first segment of the story
	segment, err := helpers.GenerateInitStory(initialPrompt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate story segment."})
		return
	}

	// Return the first segment of the story
	ctx.JSON(http.StatusOK, gin.H{"story_segment": segment[0].Content.Parts})
}
