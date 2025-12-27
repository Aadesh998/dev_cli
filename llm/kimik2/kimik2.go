package kimik2

import (
	"cli/chat"
	"cli/config"
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Kimik2Provider struct{}

type Kimik2Response struct {
	Reply string `json:"reply"`
}

func (k Kimik2Provider) ChatProcess(input string) (chat.GenericResponse, error) {
	client := openai.NewClient(
		option.WithAPIKey(config.LlmClient.ModelAPI),
		option.WithBaseURL("https://api.moonshot.ai/v1"),
	)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are Kimi, a helpful AI assistant."),
		openai.UserMessage(input),
	}

	chatCompletion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model:    openai.ChatModel("kimi-k2-turbo-preview"),
		Messages: messages,
	})

	if err != nil {
		return chat.GenericResponse{}, fmt.Errorf("Kimi request failed: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return chat.GenericResponse{}, fmt.Errorf("empty response from Kimi")
	}

	content := strings.TrimSpace(chatCompletion.Choices[0].Message.Content)

	return chat.GenericResponse{
		Text: content,
	}, nil
}
