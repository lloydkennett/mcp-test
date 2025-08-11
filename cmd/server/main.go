package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultJiraBaseURL = "https://your-domain.atlassian.net"
	defaultJiraUser    = "your-email@example.com"
	defaultJiraToken   = "your-api-token"
)

type JiraStatusRequest struct {
	TicketID string `json:"ticketId"`
	Status   string `json:"status"`
}

type JiraCommentRequest struct {
	TicketID string `json:"ticketId"`
	Comment  string `json:"comment"`
}

type JiraResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type StyleGuideRequest struct {
	Language string `json:"language"`
}

type StyleGuideResponse struct {
	Success bool   `json:"success"`
	Content string `json:"content"`
	Message string `json:"message"`
}

var (
	jiraClient  *JiraClient
	styleClient *StyleClient
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setJiraStatus(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[JiraStatusRequest]) (*mcp.CallToolResultFor[JiraResponse], error) {
	ticketID := params.Arguments.TicketID
	desiredStatus := params.Arguments.Status

	if err := jiraClient.SetTicketStatus(ctx, ticketID, desiredStatus); err != nil {
		return nil, fmt.Errorf("failed to set Jira status: %w", err)
	}

	response := JiraResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully moved ticket %s to %s", ticketID, desiredStatus),
	}

	jsonResponse, _ := json.Marshal(response)
	return &mcp.CallToolResultFor[JiraResponse]{
		Content: []mcp.Content{&mcp.TextContent{Text: string(jsonResponse)}},
		IsError: false,
	}, nil
}

func addJiraComment(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[JiraCommentRequest]) (*mcp.CallToolResultFor[JiraResponse], error) {
	ticketID := params.Arguments.TicketID
	comment := params.Arguments.Comment

	if err := jiraClient.AddComment(ctx, ticketID, comment); err != nil {
		return nil, fmt.Errorf("failed to add Jira comment: %w", err)
	}

	response := JiraResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully added comment to ticket %s", ticketID),
	}

	jsonResponse, _ := json.Marshal(response)
	return &mcp.CallToolResultFor[JiraResponse]{
		Content: []mcp.Content{&mcp.TextContent{Text: string(jsonResponse)}},
		IsError: false,
	}, nil
}

func getStyleGuide(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[StyleGuideRequest]) (*mcp.CallToolResultFor[StyleGuideResponse], error) {
	language := params.Arguments.Language

	guide, err := styleClient.GetStyleGuide(ctx, language)
	if err != nil {
		return nil, fmt.Errorf("failed to get style guide: %w", err)
	}

	response := StyleGuideResponse{
		Success: true,
		Content: fmt.Sprint(guide),
		Message: fmt.Sprintf("Retrieved %s style guide", language),
	}

	jsonResponse, _ := json.Marshal(response)
	slog.Info(string(jsonResponse))
	return &mcp.CallToolResultFor[StyleGuideResponse]{
		Content: []mcp.Content{&mcp.TextContent{Text: string(jsonResponse)}},
		IsError: false,
	}, nil
}

func main() {
	jiraBaseURL := getEnvOrDefault("JIRA_BASE_URL", defaultJiraBaseURL)
	jiraUser := getEnvOrDefault("JIRA_USER", defaultJiraUser)
	jiraToken := getEnvOrDefault("JIRA_TOKEN", defaultJiraToken)

	jiraClient = NewJiraClient(jiraBaseURL, jiraUser, jiraToken)
	styleClient = NewStyleClient()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "devtools",
		Version: "v1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "set-status",
		Description: "Set a Jira ticket to a given status (e.g., In Progress, Done, To Do)",
	}, setJiraStatus)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add-comment",
		Description: "Add a comment to a Jira ticket",
	}, addJiraComment)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "style-guide",
		Description: "Get Google's Style Guide information for different programming languages",
	}, getStyleGuide)

	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}
