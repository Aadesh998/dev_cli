package kimik2

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cli/internal/app/chat"
	"cli/internal/utils"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Kimik2Provider struct {
	client  openai.Client
	model   string
	timeout time.Duration
}

func NewKimik2Provider(apiKey, baseURL, model string) (chat.ChatProvider, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("kimi API key is missing")
	}

	if strings.TrimSpace(baseURL) == "" {
		baseURL = utils.Kimik2DefaultBaseURL
	}

	if strings.TrimSpace(model) == "" {
		model = utils.Kimik2DefaultModel
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	return &Kimik2Provider{
		client:  client,
		model:   model,
		timeout: utils.DefaultChatTimeout,
	}, nil
}

func (k *Kimik2Provider) ChatProcess(input string) (chat.GenericResponse, error) {
	if strings.TrimSpace(input) == "" {
		return chat.GenericResponse{}, errors.New("empty prompt")
	}

	ctx, cancel := context.WithTimeout(context.Background(), k.timeout)
	defer cancel()

	resp, err := k.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: openai.ChatModel(k.model),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are Kimi, a helpful AI assistant."),
				openai.UserMessage(input),
			},
			Temperature: openai.Float(0.7),
			MaxTokens:   openai.Int(2048),
		},
	)

	if err != nil {
		return chat.GenericResponse{}, fmt.Errorf("kimi completion error: %w", err)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return chat.GenericResponse{}, errors.New("empty completion response")
	}

	return chat.GenericResponse{
		Text: strings.TrimSpace(resp.Choices[0].Message.Content),
	}, nil
}
