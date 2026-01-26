package openai

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

type OpenAIProvider struct {
	client  openai.Client
	model   string
	timeout time.Duration
}

func NewOpenAIProvider(apiKey, baseURL, model string) (chat.ChatProvider, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("OpenAI API key is missing")
	}

	if strings.TrimSpace(baseURL) == "" {
		baseURL = utils.OpenAIDefaultBaseURL
	}

	if strings.TrimSpace(model) == "" {
		model = utils.OpenAIDefaultModel
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	return &OpenAIProvider{
		client:  client,
		model:   model,
		timeout: utils.DefaultChatTimeout,
	}, nil
}

func (o *OpenAIProvider) ChatProcess(input string) (chat.GenericResponse, error) {
	if strings.TrimSpace(input) == "" {
		return chat.GenericResponse{}, errors.New("empty prompt")
	}

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	resp, err := o.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: openai.ChatModel(o.model),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a helpful AI assistant."),
				openai.UserMessage(input),
			},
			Temperature: openai.Float(0.7),
			MaxTokens:   openai.Int(2048),
		},
	)

	if err != nil {
		return chat.GenericResponse{}, fmt.Errorf("openai completion error: %w", err)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return chat.GenericResponse{}, errors.New("empty completion response")
	}

	return chat.GenericResponse{
		Text: strings.TrimSpace(resp.Choices[0].Message.Content),
	}, nil
}
