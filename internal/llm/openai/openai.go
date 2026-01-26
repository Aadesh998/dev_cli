package openai

import (
	"cli/internal/app/chat"
	"cli/internal/config"
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenaiProvider struct{}

type OpenaiResponse struct {
	Reply string `json:"reply"`
}

// TODO: The ChatProcess function for OpenAI is using the Kimi model.
// This should be corrected to use an OpenAI model.
func (o OpenaiProvider) ChatProcess(input string) (chat.GenericResponse, error) {
	client := openai.NewClient(
		option.WithAPIKey(config.LlmClient.ModelAPI),
	)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful AI assistant."),
		openai.UserMessage(input),
	}

	chatCompletion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	})

	if err != nil {
		return chat.GenericResponse{}, fmt.Errorf("OpenAI request failed: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return chat.GenericResponse{}, fmt.Errorf("empty response from OpenAI")
	}

	content := strings.TrimSpace(chatCompletion.Choices[0].Message.Content)

	return chat.GenericResponse{
		Text: content,
	}, nil
}
