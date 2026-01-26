package utils

import "time"

type ModelType string

const (
	ModelClaude ModelType = "claude"
	ModelKimiK2 ModelType = "kimi-k2"
	ModelOpenAI ModelType = "openai"
	ModelGLM    ModelType = "glm"
)

const (
	// Default Timeouts
	DefaultChatTimeout = 30 * time.Second

	// Claude Constants
	ClaudeDefaultModel = "claude-3-5-sonnet-20240620"

	// GLM Constants
	GLMDefaultBaseURL = "https://open.bigmodel.cn/api/paas/v4"
	GLMDefaultModel   = "glm-4.7"

	// KimiK2 Constants
	Kimik2DefaultBaseURL = "https://api.moonshot.ai/v1"
	Kimik2DefaultModel   = "kimi-k2-turbo-preview"

	// OpenAI Constants
	OpenAIDefaultBaseURL = "https://api.openai.com/v1"
	OpenAIDefaultModel   = "gpt-5.2-turbo"
)
