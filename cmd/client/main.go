package main

import (
	"context"
	"log"
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
		Name:      "style-guide",
		Arguments: map[string]any{"language": "go"},
	}
	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatal("tool failed")
	}
	for _, c := range res.Content {
		log.Print(c.(*mcp.TextContent).Text)
	}
}
