package helpers

import (
	"context"
	"fmt"
	"os"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GenerateChatResponse continues the chat based on the user's message and chat history
func GenerateChatResponse(chatHistory string, userMessage string) ([]*genai.Candidate, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Update the prompt with the user's message and previous chat history
	updatedPrompt := fmt.Sprintf("%s\nUser said: %s\nContinue the conversation.", chatHistory, userMessage)

	// Generate content based on the updated prompt
	resp, err := client.GenerativeModel("gemini-1.5-flash").GenerateContent(ctx, genai.Text(updatedPrompt))
	if err != nil {
		return nil, fmt.Errorf("error generating chat content: %v", err)
	}

	return resp.Candidates, nil
}
