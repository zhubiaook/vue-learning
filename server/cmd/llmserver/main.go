package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	calculatorTool = mcp.NewTool("calculator",
		mcp.WithDescription("a tool for arithmetic operations, including addition, subtraction, multiplication, and division"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Enum("add", "sub", "mul", "div"),
			mcp.Description("The operation to perform"),
		),
		mcp.WithNumber("num1",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("num2",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)
	helloWorldTool = mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// add calculator tool
	s.AddTool(calculatorTool, calculatorHandler)
	// add hello world tool
	s.AddTool(helloWorldTool, helloHandler)

	svr := server.NewStreamableHTTPServer(s, server.WithEndpointPath("/mcp"))
	log.Fatal(svr.Start(":8001"))
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}

func calculatorHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	operation, err := request.RequireString("operation")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	num1, err := request.RequireFloat("num1")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	num2, err := request.RequireFloat("num2")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var result float64
	switch operation {
	case "add":
		result = num1 + num2
	case "sub":
		result = num1 - num2
	case "mul":
		result = num1 * num2
	case "div":
		if num2 == 0 {
			return mcp.NewToolResultError("Division by zero"), nil
		}
		result = num1 / num2
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid operation: %s", operation)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("The result is %f", result)), nil
}
