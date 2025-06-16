package server

import (
	"context"
	"demo/llmapp"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func StartServer(port string) error {
	r := gin.Default()
	// allow cors requests from all origins
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	r.POST("/query", handleQuery)
	r.GET("/chat", handleChat)
	return r.Run(port)
}

type QueryRequest struct {
	Message string `json:"message"`
}

func handleQuery(c *gin.Context) {
	app, err := llmapp.NewLLMApplication()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := strings.TrimSpace(req.Message)

	response, err := app.ProcessUserQuery(context.Background(), query)
	if err != nil {
		log.Println("error processing user query", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": "", "data": response})
}

// SSE event stream
func handleChat(c *gin.Context) {
	message := c.Query("message")

	client := openai.NewClient(
		option.WithAPIKey("sk-afcb4a3c1eff44ffa8bb313900c977c1"),
		option.WithBaseURL("https://api.deepseek.com"),
	)

	stream := client.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
	})

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	for stream.Next() {
		chunk := stream.Current()
		c.SSEvent("chat", gin.H{
			"content": chunk.Choices[0].Delta.Content,
		})
		c.Writer.Flush()
	}
}
