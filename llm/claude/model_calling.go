package claude

import (
	"cli/config"
	"context"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func ClaudeModelReply(input string) string {
	client := anthropic.NewClient(
		option.WithAPIKey(config.LlmClient.ModelAPI),
	)

	stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(input)),
		},
	})

	message := anthropic.Message{}
	var completeAIMesg strings.Builder
	for stream.Next() {
		event := stream.Current()
		err := message.Accumulate(event)
		if err != nil {
			panic(err)
		}

		switch eventVariant := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			switch deltaVariant := eventVariant.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				completeAIMesg.WriteString(deltaVariant.Text)
			}

		}
	}

	if stream.Err() != nil {
		panic(stream.Err())
	}
	return completeAIMesg.String()
}
