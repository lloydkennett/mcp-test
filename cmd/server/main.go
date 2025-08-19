package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultGitlabBaseURL = "https://gitlab.com"
	defaultGitlabToken   = ""
)

type GitlabProjectArgs struct {
	ProjectIDOrPath string `json:"projectIdOrPath"`
}

var (
	gitlabClient *GitlabClient
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getGitlabProject(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GitlabProjectArgs]) (*mcp.CallToolResultFor[any], error) {
	project, err := gitlabClient.GetProject(ctx, params.Arguments.ProjectIDOrPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get GitLab project: %w", err)
	}

	content, err := json.Marshal(project)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GitLab project: %w", err)
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: string(content)}},
	}, nil
}

func main() {
	gitlabBaseURL := getEnvOrDefault("GITLAB_BASE_URL", defaultGitlabBaseURL)
	gitlabToken := getEnvOrDefault("GITLAB_TOKEN", defaultGitlabToken)

	gitlabClient = NewGitlabClient(gitlabBaseURL, gitlabToken)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "devtools",
		Version: "v1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-gitlab-project",
		Description: "Get GitLab project details by project ID or URL-encoded project path",
	}, getGitlabProject)

	log.Println("Starting devtools MCP...")
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}
