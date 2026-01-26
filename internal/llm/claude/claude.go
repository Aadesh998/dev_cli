package claude

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cli/internal/app/chat"
	"cli/internal/utils"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type ClaudeProvider struct {
	client  anthropic.Client
	model   string
	timeout time.Duration
}

func NewClaudeProvider(apiKey, model string) (chat.ChatProvider, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("Claude API key is missing")
	}

	if strings.TrimSpace(model) == "" {
		model = utils.ClaudeDefaultModel
	}

	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &ClaudeProvider{
		client:  client,
		model:   model,
		timeout: utils.DefaultChatTimeout,
	}, nil
}

func (c *ClaudeProvider) ChatProcess(input string) (chat.GenericResponse, error) {
	if strings.TrimSpace(input) == "" {
		return chat.GenericResponse{}, errors.New("empty prompt")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	stream := c.client.Messages.NewStreaming(
		ctx,
		anthropic.MessageNewParams{
			Model:     anthropic.Model(c.model),
			MaxTokens: 4000,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(
					anthropic.NewTextBlock(input),
				),
			},
		},
	)

	var (
		buf strings.Builder
		msg anthropic.Message
	)

	for stream.Next() {
		event := stream.Current()

		if err := msg.Accumulate(event); err != nil {
			return chat.GenericResponse{}, fmt.Errorf("claude stream accumulate error: %w", err)
		}

		if delta, ok := event.AsAny().(anthropic.ContentBlockDeltaEvent); ok {
			if text, ok := delta.Delta.AsAny().(anthropic.TextDelta); ok {
				buf.WriteString(text.Text)
			}
		}
	}

	if err := stream.Err(); err != nil {
		return chat.GenericResponse{}, fmt.Errorf("claude stream error: %w", err)
	}

	out := strings.TrimSpace(buf.String())
	if out == "" {
		return chat.GenericResponse{}, errors.New("empty completion response")
	}

	return chat.GenericResponse{
		Text: out,
	}, nil
}
