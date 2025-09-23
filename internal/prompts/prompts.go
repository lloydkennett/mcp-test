package prompts

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Register(server *mcp.Server) {
	server.AddPrompt(&mcp.Prompt{
		Name:        "hello",
		Description: "Greet a provided name.",
		Arguments: []*mcp.PromptArgument{
			{Name: "name", Description: "Name to greet", Required: true},
		},
	}, helloHandler)
}

func helloHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	name := "world"
	if req != nil && req.Params != nil && req.Params.Arguments != nil {
		if v, ok := req.Params.Arguments["name"]; ok && v != "" {
			name = v
		}
	}
	content := &mcp.TextContent{Text: "Hello, " + name + "!"}
	return &mcp.GetPromptResult{Messages: []*mcp.PromptMessage{{Role: "user", Content: content}}}, nil
}
