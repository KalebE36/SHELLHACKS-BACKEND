
package chat

import (
	"net/http"
	"fmt"
	"SHELLHACKS-BACKEND/helpers"

	"github.com/gin-gonic/gin"
)

// StartChatHandler handles the request to start the chat
func StartChatHandler(ctx *gin.Context) {
    initialPrompt := "You are a conversational AI. Start a friendly conversation with the user. Make it engaging and ask a few questions to get to know the user better." // Initial chat prompt

    // Log the initial prompt
    fmt.Println("Initial prompt:", initialPrompt)

    // Call Gemini API to generate the first chat response
    chatSegment, err := helpers.GenerateInitChat(initialPrompt)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start the chat."})
        return
    }

    // Log the chat segment
    fmt.Println("Chat segment:", chatSegment)

    // Return the first chat segment
    ctx.JSON(http.StatusOK, gin.H{"chat_segment": chatSegment[0].Content.Parts})
}