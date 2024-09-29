package helpers

import (
	"context"
	"fmt"
	"os"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GenerateInitChat generates the initial chat response from the AI
func GenerateInitChat(prompt string) ([]*genai.Candidate, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Generate content based on the prompt
	resp, err := client.GenerativeModel("gemini-1.5-flash").GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("error generating chat content: %v", err)
	}

	return resp.Candidates, nil
}