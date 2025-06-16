package mcpclient

import (
	"context"
	"fmt"
	"sync"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	TransportHTTP = "http"
	TransportSSE  = "sse"
)

type MCPServerParam struct {
	Name              string
	Address           string
	TransportProtocol string
}

type MultiServerClient struct {
	clients map[string]*client.Client
	mutex   sync.RWMutex
}

func NewMultiServerClient() *MultiServerClient {
	return &MultiServerClient{
		clients: make(map[string]*client.Client),
	}
}

func (msc *MultiServerClient) addServer(param MCPServerParam) error {
	if _, exists := msc.clients[param.Name]; exists {
		return fmt.Errorf("mcp server %s already exists", param.Name)
	}

	msc.mutex.Lock()
	defer msc.mutex.Unlock()

	var c *client.Client

	switch param.TransportProtocol {
	case TransportHTTP:
		transport, err := transport.NewStreamableHTTP(param.Address)
		if err != nil {
			return fmt.Errorf("failed to create HTTP transport: %w", err)
		}
		c = client.NewClient(transport)
	case TransportSSE:
		transport, err := transport.NewSSE(param.Address)
		if err != nil {
			return fmt.Errorf("failed to create SSE transport: %w", err)
		}
		c = client.NewClient(transport)
	default:
		return fmt.Errorf("unsupported transport protocol: %s", param.TransportProtocol)
	}
	c.Start(context.Background())

	ctx := context.Background()
	if _, err := c.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: "1.0.0",
			Capabilities:    mcp.ClientCapabilities{},
			ClientInfo:      mcp.Implementation{},
		},
	}); err != nil {
		return fmt.Errorf("failed to initialize client for %s: %w", param.Name, err)
	}

	msc.clients[param.Name] = c
	return nil
}

func (msc *MultiServerClient) AddServers(params []MCPServerParam) error {
	for _, param := range params {
		if err := msc.addServer(param); err != nil {
			return err
		}
	}
	return nil
}

func (msc *MultiServerClient) CallTool(ctx context.Context, serverName, toolName string, args map[string]any) (*mcp.CallToolResult, error) {
	msc.mutex.RLock()
	c, exists := msc.clients[serverName]
	msc.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("server not found: %s", serverName)
	}

	return c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: args,
		},
	})
}

func (msc *MultiServerClient) AllTools(ctx context.Context) (map[string][]mcp.Tool, error) {
	msc.mutex.RLock()
	defer msc.mutex.RUnlock()

	allTools := make(map[string][]mcp.Tool)

	for serverName, c := range msc.clients {
		tools, err := c.ListTools(ctx, mcp.ListToolsRequest{})
		if err != nil {
			return nil, fmt.Errorf("failed to get tools from %s: %w", serverName, err)
		}
		allTools[serverName] = tools.Tools
	}

	return allTools, nil
}
