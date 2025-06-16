package llmclient

import (
	"context"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type LLMClient struct {
	client         openai.Client
	model          string
	tools          []openai.ChatCompletionToolParam
	messages       []openai.ChatCompletionMessageParamUnion
	requestTimeout time.Duration
}

func NewLLMClient(baseURL, apiKey, model string) *LLMClient {
	return &LLMClient{
		client: openai.NewClient(
			option.WithAPIKey(apiKey),
			option.WithBaseURL(baseURL),
		),
		model:          model,
		requestTimeout: 600 * time.Second,
	}
}

func (lc *LLMClient) AddTools(tools ...openai.ChatCompletionToolParam) {
	lc.tools = append(lc.tools, tools...)
}

func (lc *LLMClient) AddMessage(message openai.ChatCompletionMessageParamUnion) {
	lc.messages = append(lc.messages, message)
}

func (lc *LLMClient) Chat() (openai.ChatCompletionMessage, error) {
	params := openai.ChatCompletionNewParams{
		Messages: lc.messages,
		Tools:    lc.tools,
		Model:    lc.model,
	}

	ctx, cancel := context.WithTimeout(context.Background(), lc.requestTimeout)
	defer cancel()
	completion, err := lc.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}
	return completion.Choices[0].Message, nil
}

// Reset starts a new conversation by clearing all chat messages and returns the client for method chaining
func (lc *LLMClient) Reset() *LLMClient {
	lc.messages = nil
	return lc
}
