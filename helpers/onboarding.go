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
func Onboarding(scale int) ([]*genai.Candidate, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer client.Close()

	var prompt string
	// Customize the prompt based on the number
	if scale >= 0 && scale <= 24 {
		prompt = "write a small paragraph using english words for someone who barely speaks english. This person would be classified as a beginner in english. Dont make it so that you are doing things like [your name] or other. You are supposed to be a learning tool for them and introducing them to very simple words. I want you to make this about 2-3 sentences. I want you to come up with something relatively different each time but along the same level of difficulty."
	} else if scale >= 25 && scale <= 49 {
		prompt = "Write a paragraph using who would consider themselves a 2/4 when it comes to speaking english. Your paragraph should be around 3-4 sentences and should contain intermediate level words and grammar. Do not try to kill the user because you are supposed to be utilized as a learning platform. The user is going to use your reply as a way to learn new words and reinforce words they already know. Also, I want you to reply differently each time. You can make up stories or do something cool dont reply with super similar answers each time."
	} else if scale >= 50 && scale <= 74 {
		prompt = "Write a paragraph using who would consider themselves a 3/4 when it comes to speaking english. Your paragraph should be around 3-4 sentences and should contain advanced level words and grammar. You should try to use words that would be commonly used in conversation and would be considered advanced but not extreme, you are supposed to be utilized as a learning platform. The user is going to use your reply as a way to learn new words and reinforce words they already know. Also, I want you to reply differently each time. You can make up stories or do something cool dont reply with super similar answers each time. Do not talk as if you are the first person but tell a story or maybe pull something from english literature."
	} else if scale >= 75 && scale <= 100 {
		prompt = "Write a paragraph using who would consider themselves a 4/4 when it comes to speaking english. Your paragraph should be around 3-4 sentences and should contain highly advanced / extreme level words and grammar. You should try to use words that would probably not be commonly used in conversation and would be considered extreme, you are supposed to be utilized as a learning platform. The user is going to use your reply as a way to learn new words and reinforce words they already know. Also, I want you to reply differently each time. You can make up stories or do something cool dont reply with super similar answers each time. Do not talk as if you are the first person but tell a story or maybe pull something from english literature."
	}

	// Generate content based on the prompt
	resp, err := client.GenerativeModel("gemini-1.5-flash").GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("error generating content: %v", err)
	}

	return resp.Candidates, nil
}
