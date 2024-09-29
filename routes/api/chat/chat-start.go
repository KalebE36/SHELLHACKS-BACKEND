package chat

import (
	"net/http"

	"SHELLHACKS-BACKEND/helpers"

	"github.com/gin-gonic/gin"
)

// StartChatHandler handles the request to start the chat
func StartChatHandler(ctx *gin.Context) {
    initialPrompt := "Give a concise 2-sentence answer about the topic prompt given to you, which will be about language learning for spanish." 

    // Since there is no user response at the start of the chat, we pass an empty string as the second argument
    segment, err := helpers.GenerateInitChat(initialPrompt, "")
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start the chat."})
        return
    }

    // Return the first chat segment
    ctx.JSON(http.StatusOK, gin.H{"chat_segment": segment[0].Content.Parts})
}
