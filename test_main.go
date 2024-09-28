// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/google/generative-ai-go/genai"
// 	"github.com/joho/godotenv"
// 	"google.golang.org/api/option"
// )

// func main() {
// 	// Load environment variables from .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Set up the context
// 	ctx := context.Background()

// 	// Get the API key from the environment
// 	apiKey := os.Getenv("GEMINI_API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("API key is missing. Please set the GEMINI_API_KEY environment variable in the .env file.")
// 	}

// 	// Set up the Google Generative AI client
// 	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Close()

// 	// Use the generative model
// 	model := client.GenerativeModel("gemini-1.5-flash")

// 	// Generate content using the model
// 	resp, err := model.GenerateContent(ctx, genai.Text("Write a story about a magic backpack."))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Dereference and print the content of each candidate
// 	for _, candidate := range resp.Candidates {
// 		if candidate.Content != nil {
// 			fmt.Println("Generated Text:", *candidate.Content) // Dereference the Content pointer to access the string
// 		} else {
// 			fmt.Println("No content generated for this candidate.")
// 		}
// 	}
// }
