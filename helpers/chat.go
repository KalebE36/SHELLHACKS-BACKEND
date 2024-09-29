package helpers

import (
	"context"
	"fmt"
	"os"
	
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateInitChat(currentStory string, userResponse string) ([]*genai.Candidate, error) {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
    if err != nil {
        return nil, err
    }
    defer client.Close()
    
    updatedPrompt := fmt.Sprintf("%s\nUser responded: %s\nContinue the story.", currentStory, userResponse)

    // Generate content based on the updated prompt
    resp, err := client.GenerativeModel("gemini-1.5-flash").GenerateContent(ctx, genai.Text(updatedPrompt))
    if err != nil {
        return nil, fmt.Errorf("error generating content: %v", err)
    }

    return resp.Candidates, nil
}
