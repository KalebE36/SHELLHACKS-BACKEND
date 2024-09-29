package helpers

import (
	"SHELLHACKS-BACKEND/models"
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func TranslateCard(card *models.Flashcard, targetLanguage string) ([]*genai.Candidate, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Prepare the prompt for translation
	prompt := fmt.Sprintf("Translate the following word to %s: %s. Keep your response to only one word or however many words the translation is. In addition to this, I want you to add a hyphen after the word translation and include a definition also in Spanish.", targetLanguage, card.Word)

	// Call the Gemini API with the prompt
	resp, err := client.GenerativeModel("gemini-1.5-flash").GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("error generating translation: %v", err)
	}

	// Return the translated text from the response candidates
	if len(resp.Candidates) > 0 {
		return resp.Candidates, nil
	}

	return nil, err
}
