package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"SHELLHACKS-BACKEND/helpers"

	"github.com/gin-gonic/gin"
)

// GenerateParagraphsHandler handles the request to generate paragraphs based on the provided number
func GenerateParagraphsHandler(ctx *gin.Context) {
	// Get the "number" query parameter and convert it to an integer
	numParam := ctx.DefaultQuery("number", "1")
	numParagraphs, err := strconv.Atoi(numParam)
	if err != nil || numParagraphs <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid number parameter. It must be a positive integer.",
		})
		return
	}

	// Call the GenerateParagraphs function to get the paragraphs
	paragraphs, err := helpers.GenerateParagraphs(numParagraphs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate paragraphs",
		})
		return
	} else {
		fmt.Println("NOT IN HERE.")
	}

	// Return the paragraphs as a JSON response
	ctx.JSON(http.StatusOK, gin.H{
		"paragraphs": paragraphs,
	})
}
