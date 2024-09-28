package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateContentHandler(c *gin.Context) {
	ctx := context.Background()

	// Set up the Google Generative AI client
	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Gemini client"})
		return
	}
	defer client.Close()

	// Generate content using the Gemini model
	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text("Write a story about a magic backpack."))
	if err != nil {
		log.Fatalf("Failed to generate content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate content"})
		return
	}

	// Collect the generated paragraphs
	var paragraphs []string
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil {
			fmt.Printf("Content: %+v\n", *candidate.Content) // Print the structure of Content to identify the correct field
		}
	}

	// Return the paragraphs in the response (this part will be filled once the correct field is identified)
	c.JSON(http.StatusOK, gin.H{"paragraphs": paragraphs})
}
