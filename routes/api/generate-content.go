package api

import (
	"net/http"

	"SHELLHACKS-BACKEND/helpers"

	"github.com/gin-gonic/gin"
)

// GenerateParagraphsHandler handles the request to generate paragraphs based on the provided number
func GenerateParagraphsHandler(ctx *gin.Context) {
	var requestBody struct {
		Number int `json:"number"`
	}

	// Bind the JSON from the request body
	if err := ctx.ShouldBindJSON(&requestBody); err != nil || requestBody.Number <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid number parameter. It must be a positive integer.",
		})
		return
	}

	// Call the GenerateParagraphs function to get the paragraphs
	paragraphs, err := helpers.GenerateParagraphs(requestBody.Number)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate paragraphs",
		})
		return
	}

	// Return the paragraphs as a JSON response
	ctx.JSON(http.StatusOK, gin.H{
		"paragraphs": paragraphs,
	})

}
