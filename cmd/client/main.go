package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)
	transport := mcp.NewCommandTransport(exec.Command("/home/xanq/mcp-test/app/server"))

	session, err := client.Connect(ctx, transport)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	params := &mcp.CallToolParams{
		Name:      "get-gitlab-project",
		Arguments: map[string]any{"projectIdOrPath": "123"},
	}

	res, err := session.CallTool(ctx, params)
	if err != nil {
		slog.Error("CallTool failed", "error", err)
	}

	slog.Info("Tool response", "res", res.Content)

	for i, c := range res.Content {
		slog.Info("Content item", "index", i, "type", fmt.Sprintf("%T", c), "content", c)
		if textContent, ok := c.(*mcp.TextContent); ok {
			log.Printf("Content[%d] Text: %s", i, textContent.Text)
		} else {
			log.Printf("Content[%d] is not TextContent: %v", i, c)
		}
	}
}
