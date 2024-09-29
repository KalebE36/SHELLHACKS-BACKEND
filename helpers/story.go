package helpers

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GenerateParagraphsHandler handles the request to generate paragraphs based on the provided number
func GenerateStoryHandler(currentStory string, userResponse string) ([]*genai.Candidate, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}
	defer client.Close()
	updatedPrompt := fmt.Sprintf("Do not talk in the first person, this is a story for the end user. You have to act like they are the main character which is the person you are talking to. Make the continued story before the options no more than 3-4 sentences. Just give like a couple of sentences and give 3 options for how they should respond. %s\nUser responded: %s\nContinue the story. ", userResponse, currentStory)

	// Generate content based on the updated prompt
	resp, err := client.GenerativeModel("gemini-1.5-flash").GenerateContent(ctx, genai.Text(updatedPrompt))
	if err != nil {
		return nil, fmt.Errorf("error generating content: %v", err)
	}

	return resp.Candidates, nil
}
