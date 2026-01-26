package glm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cli/internal/app/chat"
	"cli/internal/utils"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type GLMProvider struct {
	client  openai.Client
	model   string
	timeout time.Duration
}

func NewGLMProvider(apiKey, baseURL, model string) (chat.ChatProvider, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("GLM API key is missing")
	}

	if strings.TrimSpace(baseURL) == "" {
		baseURL = utils.GLMDefaultBaseURL
	}

	if strings.TrimSpace(model) == "" {
		model = utils.GLMDefaultModel
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	return &GLMProvider{
		client:  client,
		model:   model,
		timeout: utils.DefaultChatTimeout,
	}, nil
}

func (g *GLMProvider) ChatProcess(message string) (chat.GenericResponse, error) {
	if strings.TrimSpace(message) == "" {
		return chat.GenericResponse{}, errors.New("empty prompt")
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	resp, err := g.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: openai.ChatModel(g.model),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are GLM, a helpful AI assistant."),
				openai.UserMessage(message),
			},
			Temperature: openai.Float(0.7),
			MaxTokens:   openai.Int(2048),
		},
	)

	if err != nil {
		return chat.GenericResponse{}, fmt.Errorf("glm completion error: %w", err)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return chat.GenericResponse{}, errors.New("empty completion response")
	}

	return chat.GenericResponse{
		Text: strings.TrimSpace(resp.Choices[0].Message.Content),
	}, nil
}
