package helpers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	
)

// GenerateParagraphs generates a specified number of paragraphs using the Gemini API based on the input number
func GenerateParagraphs(numParagraphs int) ([]*genai.Candidate, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer client.Close()

	var prompt string
	// Customize the prompt based on the number
	switch numParagraphs {
	case 3:
		prompt = "Write three paragraphs about the impact of technology on education."
	case 4:
		prompt = "Write four paragraphs about the future of space exploration."
	case 5:
		prompt = "Write five paragraphs describing the effects of climate change."
	default:
		prompt = fmt.Sprintf("Write %d paragraphs about general knowledge.", numParagraphs)
	}

	// Generate content based on the prompt
	resp, err := client.GenerativeModel("gemini-1.5-flash").GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("error generating content: %v", err)
	}

	return resp.Candidates, nil
}
