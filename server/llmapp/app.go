package llmapp

import (
	"context"
	"demo/llmclient"
	"demo/mcpclient"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
)

type LLMApplication struct {
	mcpClient *mcpclient.MultiServerClient
	llmClient *llmclient.LLMClient
}

func NewLLMApplication() (*LLMApplication, error) {
	mcpClient := mcpclient.NewMultiServerClient()

	// Add MCP servers
	mcpServerParams := []mcpclient.MCPServerParam{
		{
			Name:              "server01",
			Address:           "http://localhost:8001/mcp",
			TransportProtocol: mcpclient.TransportHTTP,
		},
	}
	if err := mcpClient.AddServers(mcpServerParams); err != nil {
		return nil, fmt.Errorf("failed to add mcp servers: %w", err)
	}

	mcpToolsMap, err := mcpClient.AllTools(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get mcp tools: %w", err)
	}

	llmClient := llmclient.NewLLMClient(
		"https://api.deepseek.com",
		"sk-afcb4a3c1eff44ffa8bb313900c977c1",
		"deepseek-chat",
	)

	llmTools := make([]openai.ChatCompletionToolParam, 0)
	for serverName, serverTools := range mcpToolsMap {
		for _, tool := range serverTools {
			llmTools = append(llmTools, openai.ChatCompletionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        fmt.Sprintf("%s_%s", serverName, tool.Name),
					Description: openai.String(tool.Description),
					Parameters: openai.FunctionParameters{
						"type":       tool.InputSchema.Type,
						"properties": tool.InputSchema.Properties,
						"required":   tool.InputSchema.Required,
					},
				},
			})
		}
	}
	llmClient.AddTools(llmTools...)

	return &LLMApplication{
		mcpClient: mcpClient,
		llmClient: llmClient,
	}, nil
}

func (app *LLMApplication) ProcessUserQuery(ctx context.Context, query string) (string, error) {
	app.llmClient.AddMessage(openai.UserMessage(query))
	completionMessage, err := app.llmClient.Chat()
	if err != nil {
		return "", fmt.Errorf("failed to get chat completion: %w", err)
	}

	llmToolCalls := completionMessage.ToolCalls
	if len(llmToolCalls) == 0 {
		return completionMessage.Content, nil
	}

	app.llmClient.AddMessage(completionMessage.ToParam())
	for _, llmToolCall := range llmToolCalls {
		// Parse server name and tool name from the combined tool name
		parts := strings.Split(llmToolCall.Function.Name, "_")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid tool name format: %s", llmToolCall.Function.Name)
		}
		mcpServerName, mcpToolName := parts[0], parts[1]

		// Extract the arguments from the function call
		var args map[string]any
		if err := json.Unmarshal([]byte(llmToolCall.Function.Arguments), &args); err != nil {
			return "", fmt.Errorf("failed to parse tool %s on server %s arguments: %w", mcpToolName, mcpServerName, err)
		}

		// Call the mcp tool on the specific server
		mcpToolResult, err := app.mcpClient.CallTool(ctx, mcpServerName, mcpToolName, args)
		if err != nil {
			return "", fmt.Errorf("failed to call tool %s on server %s: %w", mcpToolName, mcpServerName, err)
		}

		if mcpToolResult.IsError {
			return "", fmt.Errorf("tool %s on server %s returned an error: %s", mcpToolName, mcpServerName, mcpToolResult.Content[0].(mcp.TextContent).Text)
		}

		result := mcpToolResult.Content[0].(mcp.TextContent).Text
		log.Printf("call mcp tool %s on server %s returned: %s", mcpToolName, mcpServerName, result)
		app.llmClient.AddMessage(openai.ToolMessage(result, llmToolCall.ID))
	}

	completionMessage, err = app.llmClient.Chat()
	if err != nil {
		return "", fmt.Errorf("failed to get final chat completion: %w", err)
	}

	return completionMessage.Content, nil
}
